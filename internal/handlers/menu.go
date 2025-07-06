package handlers

import (
	"net/http"
	"pos-system/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MenuHandler struct {
	db *gorm.DB
}

func NewMenuHandler(db *gorm.DB) *MenuHandler {
	return &MenuHandler{db: db}
}

// Categories
func (h *MenuHandler) GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := h.db.Preload("MenuItems").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *MenuHandler) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *MenuHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	
	var category models.Category
	if err := h.db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *MenuHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	
	if err := h.db.Delete(&models.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// Menu Items
func (h *MenuHandler) GetMenuItems(c *gin.Context) {
	categoryID := c.Query("category_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var menuItems []models.MenuItem
	var total int64

	query := h.db.Model(&models.MenuItem{}).Preload("Category")
	
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	query.Count(&total)
	if err := query.Offset(offset).Limit(limit).Find(&menuItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch menu items"})
		return
	}

	// Calculate margin for each item
	for i := range menuItems {
		if menuItems[i].Price > 0 {
			menuItems[i].Margin = ((menuItems[i].Price - menuItems[i].COGS) / menuItems[i].Price) * 100
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  menuItems,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *MenuHandler) CreateMenuItem(c *gin.Context) {
	var menuItem models.MenuItem
	if err := c.ShouldBindJSON(&menuItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate category exists
	var category models.Category
	if err := h.db.First(&category, menuItem.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
		return
	}

	if err := h.db.Create(&menuItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create menu item"})
		return
	}

	// Calculate margin
	if menuItem.Price > 0 {
		menuItem.Margin = ((menuItem.Price - menuItem.COGS) / menuItem.Price) * 100
	}

	c.JSON(http.StatusCreated, menuItem)
}

func (h *MenuHandler) GetMenuItem(c *gin.Context) {
	id := c.Param("id")
	
	var menuItem models.MenuItem
	if err := h.db.Preload("Category").First(&menuItem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu item not found"})
		return
	}

	// Calculate margin
	if menuItem.Price > 0 {
		menuItem.Margin = ((menuItem.Price - menuItem.COGS) / menuItem.Price) * 100
	}

	c.JSON(http.StatusOK, menuItem)
}

func (h *MenuHandler) UpdateMenuItem(c *gin.Context) {
	id := c.Param("id")
	
	var menuItem models.MenuItem
	if err := h.db.First(&menuItem, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu item not found"})
		return
	}

	if err := c.ShouldBindJSON(&menuItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate category exists if category_id is being updated
	if menuItem.CategoryID != 0 {
		var category models.Category
		if err := h.db.First(&category, menuItem.CategoryID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
			return
		}
	}

	if err := h.db.Save(&menuItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update menu item"})
		return
	}

	// Calculate margin
	if menuItem.Price > 0 {
		menuItem.Margin = ((menuItem.Price - menuItem.COGS) / menuItem.Price) * 100
	}

	c.JSON(http.StatusOK, menuItem)
}

func (h *MenuHandler) DeleteMenuItem(c *gin.Context) {
	id := c.Param("id")
	
	if err := h.db.Delete(&models.MenuItem{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete menu item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu item deleted successfully"})
}
