package handlers

import (
	"fmt"
	"log"
	"net/http"
	"pos-system/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	db *gorm.DB
}

type CreateTransactionRequest struct {
	CustomerName string                   `json:"customer_name"`
	Items        []TransactionItemRequest `json:"items" binding:"required"`
	Tax          float64                  `json:"tax"`
	Discount     float64                  `json:"discount"`
}

type TransactionItemRequest struct {
	MenuItemID uint                      `json:"menu_item_id" binding:"required"`
	Quantity   int                       `json:"quantity" binding:"required,min=1"`
	AddOns     []TransactionItemAddOnRequest `json:"add_ons,omitempty"`
}

type TransactionItemAddOnRequest struct {
	AddOnID  uint `json:"add_on_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required,min=1"`
}

type PayTransactionRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required"`
}

type UpdateTransactionRequest struct {
	CustomerName string  `json:"customer_name"`
	Tax          float64 `json:"tax"`
	Discount     float64 `json:"discount"`
}

type AddTransactionItemRequest struct {
	MenuItemID uint                      `json:"menu_item_id" binding:"required"`
	Quantity   int                       `json:"quantity" binding:"required,min=1"`
	AddOns     []TransactionItemAddOnRequest `json:"add_ons,omitempty"`
}

type UpdateTransactionItemRequest struct {
	Quantity int                       `json:"quantity" binding:"required,min=1"`
	AddOns   []TransactionItemAddOnRequest `json:"add_ons,omitempty"`
}

func NewTransactionHandler(db *gorm.DB) *TransactionHandler {
	return &TransactionHandler{db: db}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	
	// Start transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Generate transaction number
	transactionNo := fmt.Sprintf("TRX-%d", time.Now().Unix())

	// Create transaction
	transaction := models.Transaction{
		TransactionNo: transactionNo,
		UserID:        userID.(uint),
		CustomerName:  req.CustomerName,
		Status:        "pending",
		Tax:           req.Tax,
		Discount:      req.Discount,
	}

	var subTotal float64

	// Calculate subtotal and validate items
	for _, itemReq := range req.Items {
		var menuItem models.MenuItem
		if err := tx.First(&menuItem, itemReq.MenuItemID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Menu item %d not found", itemReq.MenuItemID)})
			return
		}

		if !menuItem.IsAvailable {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Menu item %s is not available", menuItem.Name)})
			return
		}

		itemTotal := menuItem.Price * float64(itemReq.Quantity)

		// Validate and calculate add-ons
		for _, addOnReq := range itemReq.AddOns {
			var addOn models.AddOn
			if err := tx.First(&addOn, addOnReq.AddOnID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Add-on %d not found", addOnReq.AddOnID)})
				return
			}

			if !addOn.IsAvailable {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Add-on %s is not available", addOn.Name)})
				return
			}

			itemTotal += addOn.Price * float64(addOnReq.Quantity) * float64(itemReq.Quantity)
		}

		subTotal += itemTotal
	}

	transaction.SubTotal = subTotal
	transaction.Total = subTotal + req.Tax - req.Discount

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	// Create transaction items and add-ons
	for _, itemReq := range req.Items {
		var menuItem models.MenuItem
		tx.First(&menuItem, itemReq.MenuItemID)

		totalPrice := menuItem.Price * float64(itemReq.Quantity)

		// Calculate add-ons total for this item
		var addOnsTotal float64
		for _, addOnReq := range itemReq.AddOns {
			var addOn models.AddOn
			tx.First(&addOn, addOnReq.AddOnID)
			addOnsTotal += addOn.Price * float64(addOnReq.Quantity) * float64(itemReq.Quantity)
		}

		transactionItem := models.TransactionItem{
			TransactionID: transaction.ID,
			MenuItemID:    itemReq.MenuItemID,
			Quantity:      itemReq.Quantity,
			UnitPrice:     menuItem.Price,
			TotalPrice:    totalPrice + addOnsTotal,
		}

		if err := tx.Create(&transactionItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction item"})
			return
		}

		// Create add-ons for this item
		for _, addOnReq := range itemReq.AddOns {
			var addOn models.AddOn
			tx.First(&addOn, addOnReq.AddOnID)

			transactionItemAddOn := models.TransactionItemAddOn{
				TransactionItemID: transactionItem.ID,
				AddOnID:           addOnReq.AddOnID,
				Quantity:          addOnReq.Quantity,
				UnitPrice:         addOn.Price,
				TotalPrice:        addOn.Price * float64(addOnReq.Quantity) * float64(itemReq.Quantity),
			}

			if err := tx.Create(&transactionItemAddOn).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction item add-on"})
				return
			}
		}
	}

	tx.Commit()

	// Reload with associations
	h.db.Preload("Items.MenuItem").
		Preload("Items.AddOns.AddOn").
		Preload("User").
		First(&transaction, transaction.ID)

	c.JSON(http.StatusCreated, transaction)
}

func (h *TransactionHandler) PayTransaction(c *gin.Context) {
	id := c.Param("id")
	
	var req PayTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("PayTransaction: Failed to bind JSON for transaction %s: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}
	
	log.Printf("PayTransaction: Processing payment for transaction %s with method %s", id, req.PaymentMethod)

	var transaction models.Transaction
	if err := h.db.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if transaction.Status == "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction already paid"})
		return
	}

	// Validate payment method
	var paymentMethod models.PaymentMethod
	if err := h.db.Where("code = ? AND is_active = ?", req.PaymentMethod, true).First(&paymentMethod).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method"})
		return
	}

	now := time.Now()
	transaction.Status = "paid"
	transaction.PaymentMethod = req.PaymentMethod
	transaction.PaidAt = &now

	if err := h.db.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
		return
	}

	// Reload with associations
	h.db.Preload("Items.MenuItem").
		Preload("Items.AddOns.AddOn").
		Preload("User").
		First(&transaction, transaction.ID)

	c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var transactions []models.Transaction
	var total int64

	query := h.db.Model(&models.Transaction{}).
		Preload("Items.MenuItem").
		Preload("Items.AddOns.AddOn").
		Preload("User")
	
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  transactions,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	id := c.Param("id")
	
	var transaction models.Transaction
	if err := h.db.Preload("Items.MenuItem").
		Preload("Items.AddOns.AddOn").
		Preload("User").
		First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	id := c.Param("id")
	
	var transaction models.Transaction
	if err := h.db.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	// if transaction.Status == "paid" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete paid transaction"})
	// 	return
	// }

	// Start transaction to delete all related records
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete transaction item add-ons first
	if err := tx.Where("transaction_item_id IN (SELECT id FROM transaction_items WHERE transaction_id = ?)", id).
		Delete(&models.TransactionItemAddOn{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction item add-ons"})
		return
	}

	// Delete transaction items
	if err := tx.Where("transaction_id = ?", id).Delete(&models.TransactionItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction items"})
		return
	}

	// Delete transaction
	if err := tx.Delete(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}

// UpdateTransaction updates basic transaction information (customer name, tax, discount)
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	id := c.Param("id")
	
	var transaction models.Transaction
	if err := h.db.First(&transaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if transaction.Status == "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update paid transaction"})
		return
	}

	var req UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update transaction fields
	transaction.CustomerName = req.CustomerName
	transaction.Tax = req.Tax
	transaction.Discount = req.Discount
	transaction.UpdatedAt = time.Now()

	// Recalculate total
	var total float64
	var items []models.TransactionItem
	if err := h.db.Preload("AddOns.AddOn").Where("transaction_id = ?", id).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction items"})
		return
	}

	for _, item := range items {
		var menuItem models.MenuItem
		if err := h.db.First(&menuItem, item.MenuItemID).Error; err != nil {
			continue
		}
		itemTotal := menuItem.Price * float64(item.Quantity)
		
		for _, addOn := range item.AddOns {
			// Use the stored TotalPrice which already includes menu item quantity
			itemTotal += addOn.TotalPrice
		}
		total += itemTotal
	}

	transaction.SubTotal = total
	transaction.Total = total + transaction.Tax - transaction.Discount

	if err := h.db.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// AddTransactionItem adds a new item to a pending transaction
func (h *TransactionHandler) AddTransactionItem(c *gin.Context) {
	transactionID := c.Param("id")
	
	var transaction models.Transaction
	if err := h.db.First(&transaction, transactionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if transaction.Status == "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot modify paid transaction"})
		return
	}

	var req AddTransactionItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify menu item exists
	var menuItem models.MenuItem
	if err := h.db.First(&menuItem, req.MenuItemID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Menu item not found"})
		return
	}

	// Start transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create transaction item
	transactionItem := models.TransactionItem{
		TransactionID: transaction.ID,
		MenuItemID:    req.MenuItemID,
		Quantity:      req.Quantity,
		UnitPrice:     menuItem.Price,
		TotalPrice:    menuItem.Price * float64(req.Quantity),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := tx.Create(&transactionItem).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction item"})
		return
	}

	// Create transaction item add-ons if provided
	for _, addOnReq := range req.AddOns {
		var addOn models.AddOn
		if err := tx.First(&addOn, addOnReq.AddOnID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Add-on with ID %d not found", addOnReq.AddOnID)})
			return
		}

		transactionItemAddOn := models.TransactionItemAddOn{
			TransactionItemID: transactionItem.ID,
			AddOnID:          addOnReq.AddOnID,
			Quantity:         addOnReq.Quantity,
			UnitPrice:        addOn.Price,
			TotalPrice:       addOn.Price * float64(addOnReq.Quantity) * float64(req.Quantity),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		if err := tx.Create(&transactionItemAddOn).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction item add-on"})
			return
		}
	}

	// Recalculate transaction totals
	var total float64
	var items []models.TransactionItem
	if err := tx.Preload("AddOns.AddOn").Where("transaction_id = ?", transactionID).Find(&items).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate totals"})
		return
	}

	for _, item := range items {
		var itemMenuItem models.MenuItem
		if err := tx.First(&itemMenuItem, item.MenuItemID).Error; err != nil {
			continue
		}
		itemTotal := itemMenuItem.Price * float64(item.Quantity)
		
		for _, addOn := range item.AddOns {
			itemTotal += addOn.TotalPrice
		}
		total += itemTotal
	}

	transaction.SubTotal = total
	transaction.Total = total + transaction.Tax - transaction.Discount
	transaction.UpdatedAt = time.Now()

	if err := tx.Save(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction totals"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, transactionItem)
}

// UpdateTransactionItem updates an existing transaction item
func (h *TransactionHandler) UpdateTransactionItem(c *gin.Context) {
	transactionID := c.Param("id")
	itemID := c.Param("item_id")
	
	var transaction models.Transaction
	if err := h.db.First(&transaction, transactionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if transaction.Status == "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot modify paid transaction"})
		return
	}

	var transactionItem models.TransactionItem
	if err := h.db.Where("id = ? AND transaction_id = ?", itemID, transactionID).First(&transactionItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction item not found"})
		return
	}

	var req UpdateTransactionItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update transaction item quantity
	transactionItem.Quantity = req.Quantity
	transactionItem.UpdatedAt = time.Now()

	if err := tx.Save(&transactionItem).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction item"})
		return
	}

	// Delete existing add-ons
	if err := tx.Where("transaction_item_id = ?", itemID).Delete(&models.TransactionItemAddOn{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction item add-ons"})
		return
	}

	// Create new add-ons
	for _, addOnReq := range req.AddOns {
		var addOn models.AddOn
		if err := tx.First(&addOn, addOnReq.AddOnID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Add-on with ID %d not found", addOnReq.AddOnID)})
			return
		}

		transactionItemAddOn := models.TransactionItemAddOn{
			TransactionItemID: transactionItem.ID,
			AddOnID:          addOnReq.AddOnID,
			Quantity:         addOnReq.Quantity,
			UnitPrice:        addOn.Price,
			TotalPrice:       addOn.Price * float64(addOnReq.Quantity) * float64(transactionItem.Quantity),
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		if err := tx.Create(&transactionItemAddOn).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction item add-on"})
			return
		}
	}

	// Recalculate transaction totals
	var total float64
	var items []models.TransactionItem
	if err := tx.Preload("AddOns.AddOn").Where("transaction_id = ?", transactionID).Find(&items).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate totals"})
		return
	}

	for _, item := range items {
		var itemMenuItem models.MenuItem
		if err := tx.First(&itemMenuItem, item.MenuItemID).Error; err != nil {
			continue
		}
		itemTotal := itemMenuItem.Price * float64(item.Quantity)
		
		for _, addOn := range item.AddOns {
			itemTotal += addOn.TotalPrice
		}
		total += itemTotal
	}

	transaction.SubTotal = total
	transaction.Total = total + transaction.Tax - transaction.Discount
	transaction.UpdatedAt = time.Now()

	if err := tx.Save(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction totals"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, transactionItem)
}

// DeleteTransactionItem removes an item from a pending transaction
func (h *TransactionHandler) DeleteTransactionItem(c *gin.Context) {
	transactionID := c.Param("id")
	itemID := c.Param("item_id")
	
	var transaction models.Transaction
	if err := h.db.First(&transaction, transactionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if transaction.Status == "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot modify paid transaction"})
		return
	}

	var transactionItem models.TransactionItem
	if err := h.db.Where("id = ? AND transaction_id = ?", itemID, transactionID).First(&transactionItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction item not found"})
		return
	}

	// Start transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete transaction item add-ons first
	if err := tx.Where("transaction_item_id = ?", itemID).Delete(&models.TransactionItemAddOn{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction item add-ons"})
		return
	}

	// Delete transaction item
	if err := tx.Delete(&transactionItem).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction item"})
		return
	}

	// Recalculate transaction totals
	var total float64
	var items []models.TransactionItem
	if err := tx.Preload("AddOns.AddOn").Where("transaction_id = ?", transactionID).Find(&items).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate totals"})
		return
	}

	for _, item := range items {
		var itemMenuItem models.MenuItem
		if err := tx.First(&itemMenuItem, item.MenuItemID).Error; err != nil {
			continue
		}
		itemTotal := itemMenuItem.Price * float64(item.Quantity)
		
		for _, addOn := range item.AddOns {
			itemTotal += addOn.TotalPrice
		}
		total += itemTotal
	}

	transaction.SubTotal = total
	transaction.Total = total + transaction.Tax - transaction.Discount
	transaction.UpdatedAt = time.Now()

	if err := tx.Save(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction totals"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Transaction item deleted successfully"})
}

func (h *TransactionHandler) GetPaymentMethods(c *gin.Context) {
	var paymentMethods []models.PaymentMethod
	if err := h.db.Where("is_active = ?", true).Find(&paymentMethods).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment methods"})
		return
	}

	c.JSON(http.StatusOK, paymentMethods)
}
