package handlers

import (
	"net/http"
	"pos-system/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExpenseHandler struct {
	db *gorm.DB
}

type CreateExpenseRequest struct {
	Type        string    `json:"type" binding:"required,oneof=raw_material operational"`
	Category    string    `json:"category" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Amount      float64   `json:"amount" binding:"required,gt=0"`
	Date        time.Time `json:"date" binding:"required"`
}

func NewExpenseHandler(db *gorm.DB) *ExpenseHandler {
	return &ExpenseHandler{db: db}
}

func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var req CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	expense := models.Expense{
		Type:        req.Type,
		Category:    req.Category,
		Description: req.Description,
		Amount:      req.Amount,
		Date:        req.Date,
		UserID:      userID.(uint),
	}

	if err := h.db.Create(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
		return
	}

	// Reload with user association
	h.db.Preload("User").First(&expense, expense.ID)

	c.JSON(http.StatusCreated, expense)
}

func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
	expenseType := c.Query("type")
	category := c.Query("category")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var expenses []models.Expense
	var total int64

	query := h.db.Model(&models.Expense{}).Preload("User")

	if expenseType != "" {
		query = query.Where("type = ?", expenseType)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	query.Count(&total)
	if err := query.Order("date DESC").Offset(offset).Limit(limit).Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  expenses,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *ExpenseHandler) GetExpense(c *gin.Context) {
	id := c.Param("id")

	var expense models.Expense
	if err := h.db.Preload("User").First(&expense, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	c.JSON(http.StatusOK, expense)
}

func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	id := c.Param("id")

	var expense models.Expense
	if err := h.db.First(&expense, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	var req CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense.Type = req.Type
	expense.Category = req.Category
	expense.Description = req.Description
	expense.Amount = req.Amount
	expense.Date = req.Date

	if err := h.db.Save(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
		return
	}

	// Reload with user association
	h.db.Preload("User").First(&expense, expense.ID)

	c.JSON(http.StatusOK, expense)
}

func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Delete(&models.Expense{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}

func (h *ExpenseHandler) GetExpenseSummary(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	query := h.db.Model(&models.Expense{})

	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	// Total expenses
	var totalExpenses float64
	query.Select("COALESCE(SUM(amount), 0)").Scan(&totalExpenses)

	// Expenses by type
	var expensesByType []struct {
		Type   string  `json:"type"`
		Amount float64 `json:"amount"`
	}
	query.Select("type, COALESCE(SUM(amount), 0) as amount").Group("type").Scan(&expensesByType)

	// Expenses by category
	var expensesByCategory []struct {
		Category string  `json:"category"`
		Amount   float64 `json:"amount"`
	}
	query.Select("category, COALESCE(SUM(amount), 0) as amount").Group("category").Scan(&expensesByCategory)

	c.JSON(http.StatusOK, gin.H{
		"total_expenses":        totalExpenses,
		"expenses_by_type":      expensesByType,
		"expenses_by_category":  expensesByCategory,
	})
}
