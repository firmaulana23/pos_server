package handlers

import (
	"net/http"
	"pos-system/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AddOnHandler struct {
	db *gorm.DB
}

func NewAddOnHandler(db *gorm.DB) *AddOnHandler {
	return &AddOnHandler{db: db}
}

func (h *AddOnHandler) GetAddOns(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var addOns []models.AddOn
	var total int64

	query := h.db.Model(&models.AddOn{}).Preload("MenuItem")
	
	// Filter by menu item ID if provided
	if menuItemID := c.Query("menu_item_id"); menuItemID != "" {
		if menuItemID == "global" {
			// Get global add-ons (not tied to any specific menu item)
			query = query.Where("menu_item_id IS NULL")
		} else {
			// Get add-ons for specific menu item OR global add-ons
			query = query.Where("menu_item_id = ? OR menu_item_id IS NULL", menuItemID)
		}
	}
	
	// Filter by availability if requested
	if available := c.Query("available"); available != "" {
		if available == "true" {
			query = query.Where("is_available = ?", true)
		}
	}

	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Find(&addOns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch add-ons"})
		return
	}

	// Calculate margin for each add-on
	for i := range addOns {
		if addOns[i].Price > 0 {
			addOns[i].Margin = ((addOns[i].Price - addOns[i].COGS) / addOns[i].Price) * 100
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  addOns,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *AddOnHandler) CreateAddOn(c *gin.Context) {
	var addOn models.AddOn
	if err := c.ShouldBindJSON(&addOn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&addOn).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create add-on"})
		return
	}

	// Calculate margin
	if addOn.Price > 0 {
		addOn.Margin = ((addOn.Price - addOn.COGS) / addOn.Price) * 100
	}

	c.JSON(http.StatusCreated, addOn)
}

func (h *AddOnHandler) GetAddOn(c *gin.Context) {
	id := c.Param("id")
	
	var addOn models.AddOn
	if err := h.db.Preload("MenuItem").First(&addOn, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Add-on not found"})
		return
	}

	// Calculate margin
	if addOn.Price > 0 {
		addOn.Margin = ((addOn.Price - addOn.COGS) / addOn.Price) * 100
	}

	c.JSON(http.StatusOK, addOn)
}

func (h *AddOnHandler) UpdateAddOn(c *gin.Context) {
	id := c.Param("id")
	
	var addOn models.AddOn
	if err := h.db.First(&addOn, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Add-on not found"})
		return
	}

	if err := c.ShouldBindJSON(&addOn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&addOn).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update add-on"})
		return
	}

	// Calculate margin
	if addOn.Price > 0 {
		addOn.Margin = ((addOn.Price - addOn.COGS) / addOn.Price) * 100
	}

	c.JSON(http.StatusOK, addOn)
}

func (h *AddOnHandler) DeleteAddOn(c *gin.Context) {
	id := c.Param("id")
	
	if err := h.db.Delete(&models.AddOn{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete add-on"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Add-on deleted successfully"})
}

// GetAddOnsForMenuItem gets all add-ons available for a specific menu item
// This includes both global add-ons and menu-item-specific add-ons
func (h *AddOnHandler) GetAddOnsForMenuItem(c *gin.Context) {
	menuItemID := c.Param("menu_item_id")
	
	// Check if menu item exists
	var menuItem models.MenuItem
	if err := h.db.First(&menuItem, menuItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu item not found"})
		return
	}

	var addOns []models.AddOn
	
	// Get add-ons for this specific menu item OR global add-ons (menu_item_id IS NULL)
	if err := h.db.Where("(menu_item_id = ? OR menu_item_id IS NULL) AND is_available = ?", menuItemID, true).
		Preload("MenuItem").Find(&addOns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch add-ons"})
		return
	}

	// Calculate margin for each add-on
	for i := range addOns {
		if addOns[i].Price > 0 {
			addOns[i].Margin = ((addOns[i].Price - addOns[i].COGS) / addOns[i].Price) * 100
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"menu_item": gin.H{
			"id":   menuItem.ID,
			"name": menuItem.Name,
		},
		"add_ons": addOns,
	})
}
