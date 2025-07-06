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

	query := h.db.Model(&models.AddOn{})
	
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
	if err := h.db.First(&addOn, id).Error; err != nil {
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
