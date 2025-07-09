package handlers

import (
	"net/http"
	"pos-system/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	db *gorm.DB
}

type DashboardStats struct {
	TotalSales     float64       `json:"total_sales"`
	TotalCOGS      float64       `json:"total_cogs"`
	GrossProfit    float64       `json:"gross_profit"`
	GrossMargin    float64       `json:"gross_margin_percent"`
	TotalExpenses  float64       `json:"total_expenses"`
	NetProfit      float64       `json:"net_profit"`
	TotalOrders    int64         `json:"total_orders"`
	PendingOrders  int64         `json:"pending_orders"`
	PaidOrders     int64         `json:"paid_orders"`
	TopMenuItems   []TopMenuItem `json:"top_menu_items"`
	TopAddOns      []TopAddOn    `json:"top_add_ons"`
	SalesChart     []SalesData   `json:"sales_chart"`
	ExpenseChart   []ExpenseData `json:"expense_chart"`
}

type TopMenuItem struct {
	Name         string  `json:"name"`
	TotalSold    int     `json:"total_sold"`
	TotalRevenue float64 `json:"total_revenue"`
}

type TopAddOn struct {
	Name         string  `json:"name"`
	TotalSold    int     `json:"total_sold"`
	TotalRevenue float64 `json:"total_revenue"`
}

type SalesData struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
	Orders int64   `json:"orders"`
}

type ExpenseData struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	stats := DashboardStats{}

	// Build base queries based on whether date filters are provided
	var salesQuery, expenseQuery, orderQuery *gorm.DB

	if startDate != "" && endDate != "" {
		// Use date filtering when dates are provided
		salesQuery = h.db.Model(&models.Transaction{}).
			Where("status = ? AND DATE(created_at) BETWEEN ? AND ?", "paid", startDate, endDate)
		expenseQuery = h.db.Model(&models.Expense{}).
			Where("DATE(date) BETWEEN ? AND ?", startDate, endDate)
		orderQuery = h.db.Model(&models.Transaction{}).
			Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate)
	} else {
		// Use all data when no date filters are provided
		salesQuery = h.db.Model(&models.Transaction{}).
			Where("status = ?", "paid")
		expenseQuery = h.db.Model(&models.Expense{})
		orderQuery = h.db.Model(&models.Transaction{})
	}

	// Total Sales (paid transactions only)
	salesQuery.Select("COALESCE(SUM(total), 0)").Scan(&stats.TotalSales)

	// Calculate Total COGS
	var cogsQuery *gorm.DB
	if startDate != "" && endDate != "" {
		cogsQuery = h.db.Table("transaction_items").
			Select("COALESCE(SUM(transaction_items.quantity * menu_items.cogs), 0)").
			Joins("JOIN menu_items ON transaction_items.menu_item_id = menu_items.id").
			Joins("JOIN transactions ON transaction_items.transaction_id = transactions.id").
			Where("transactions.status = ? AND DATE(transactions.created_at) BETWEEN ? AND ?", "paid", startDate, endDate)
	} else {
		cogsQuery = h.db.Table("transaction_items").
			Select("COALESCE(SUM(transaction_items.quantity * menu_items.cogs), 0)").
			Joins("JOIN menu_items ON transaction_items.menu_item_id = menu_items.id").
			Joins("JOIN transactions ON transaction_items.transaction_id = transactions.id").
			Where("transactions.status = ?", "paid")
	}
	cogsQuery.Scan(&stats.TotalCOGS)

	// Calculate Add-ons COGS
	var addOnCOGS float64
	var addOnCogsQuery *gorm.DB
	if startDate != "" && endDate != "" {
		addOnCogsQuery = h.db.Table("transaction_item_add_ons").
			Select("COALESCE(SUM(transaction_item_add_ons.quantity * add_ons.cogs), 0)").
			Joins("JOIN add_ons ON transaction_item_add_ons.add_on_id = add_ons.id").
			Joins("JOIN transaction_items ON transaction_item_add_ons.transaction_item_id = transaction_items.id").
			Joins("JOIN transactions ON transaction_items.transaction_id = transactions.id").
			Where("transactions.status = ? AND DATE(transactions.created_at) BETWEEN ? AND ?", "paid", startDate, endDate)
	} else {
		addOnCogsQuery = h.db.Table("transaction_item_add_ons").
			Select("COALESCE(SUM(transaction_item_add_ons.quantity * add_ons.cogs), 0)").
			Joins("JOIN add_ons ON transaction_item_add_ons.add_on_id = add_ons.id").
			Joins("JOIN transaction_items ON transaction_item_add_ons.transaction_item_id = transaction_items.id").
			Joins("JOIN transactions ON transaction_items.transaction_id = transactions.id").
			Where("transactions.status = ?", "paid")
	}
	addOnCogsQuery.Scan(&addOnCOGS)

	// Total COGS includes menu items and add-ons
	stats.TotalCOGS += addOnCOGS

	// Calculate Gross Profit (Sales - COGS)
	stats.GrossProfit = stats.TotalSales - stats.TotalCOGS

	// Calculate Gross Margin Percentage
	if stats.TotalSales > 0 {
		stats.GrossMargin = (stats.GrossProfit / stats.TotalSales) * 100
	} else {
		stats.GrossMargin = 0
	}

	// Total Expenses
	expenseQuery.Select("COALESCE(SUM(amount), 0)").Scan(&stats.TotalExpenses)

	// Net Profit (Gross Profit - Expenses)
	stats.NetProfit = stats.GrossProfit - stats.TotalExpenses

	// Order counts
	orderQuery.Count(&stats.TotalOrders)

	if startDate != "" && endDate != "" {
		h.db.Model(&models.Transaction{}).
			Where("status = ? AND DATE(created_at) BETWEEN ? AND ?", "pending", startDate, endDate).
			Count(&stats.PendingOrders)

		h.db.Model(&models.Transaction{}).
			Where("status = ? AND DATE(created_at) BETWEEN ? AND ?", "paid", startDate, endDate).
			Count(&stats.PaidOrders)
	} else {
		h.db.Model(&models.Transaction{}).
			Where("status = ?", "pending").
			Count(&stats.PendingOrders)

		h.db.Model(&models.Transaction{}).
			Where("status = ?", "paid").
			Count(&stats.PaidOrders)
	}

	// Top menu items
	topMenuQuery := h.db.Table("transaction_items").
		Select("menu_items.name, SUM(transaction_items.quantity) as total_sold, SUM(transaction_items.total_price) as total_revenue").
		Joins("JOIN menu_items ON transaction_items.menu_item_id = menu_items.id").
		Joins("JOIN transactions ON transaction_items.transaction_id = transactions.id")

	if startDate != "" && endDate != "" {
		topMenuQuery = topMenuQuery.Where("transactions.status = ? AND DATE(transactions.created_at) BETWEEN ? AND ?", "paid", startDate, endDate)
	} else {
		topMenuQuery = topMenuQuery.Where("transactions.status = ?", "paid")
	}

	topMenuQuery.Group("menu_items.id, menu_items.name").
		Order("total_sold DESC").
		Limit(5).
		Scan(&stats.TopMenuItems)

	// Top add-ons
	topAddOnQuery := h.db.Table("transaction_item_add_ons").
		Select("add_ons.name, SUM(transaction_item_add_ons.quantity) as total_sold, SUM(transaction_item_add_ons.total_price) as total_revenue").
		Joins("JOIN add_ons ON transaction_item_add_ons.add_on_id = add_ons.id").
		Joins("JOIN transaction_items ON transaction_item_add_ons.transaction_item_id = transaction_items.id").
		Joins("JOIN transactions ON transaction_items.transaction_id = transactions.id")

	if startDate != "" && endDate != "" {
		topAddOnQuery = topAddOnQuery.Where("transactions.status = ? AND DATE(transactions.created_at) BETWEEN ? AND ?", "paid", startDate, endDate)
	} else {
		topAddOnQuery = topAddOnQuery.Where("transactions.status = ?", "paid")
	}

	topAddOnQuery.Group("add_ons.id, add_ons.name").
		Order("total_sold DESC").
		Limit(5).
		Scan(&stats.TopAddOns)

	// Sales chart data
	if startDate != "" && endDate != "" {
		h.db.Raw(`
			SELECT 
				DATE(created_at) as date,
				COALESCE(SUM(total), 0) as amount,
				COUNT(*) as orders
			FROM transactions 
			WHERE status = 'paid' AND DATE(created_at) BETWEEN ? AND ?
			GROUP BY DATE(created_at)
			ORDER BY date DESC
		`, startDate, endDate).Scan(&stats.SalesChart)
	} else {
		h.db.Raw(`
			SELECT 
				DATE(created_at) as date,
				COALESCE(SUM(total), 0) as amount,
				COUNT(*) as orders
			FROM transactions 
			WHERE status = 'paid'
			GROUP BY DATE(created_at)
			ORDER BY date DESC
			LIMIT 30
		`).Scan(&stats.SalesChart)
	}

	// Expense chart data
	if startDate != "" && endDate != "" {
		h.db.Raw(`
			SELECT 
				DATE(date) as date,
				COALESCE(SUM(amount), 0) as amount,
				type
			FROM expenses 
			WHERE deleted_at IS NULL AND DATE(date) BETWEEN ? AND ?
			GROUP BY DATE(date), type
			ORDER BY date DESC
		`, startDate, endDate).Scan(&stats.ExpenseChart)
	} else {
		h.db.Raw(`
			SELECT 
				DATE(date) as date,
				COALESCE(SUM(amount), 0) as amount,
				type
			FROM expenses 
			WHERE deleted_at IS NULL
			GROUP BY DATE(date), type
			ORDER BY date DESC
			LIMIT 30
		`).Scan(&stats.ExpenseChart)
	}

	c.JSON(http.StatusOK, stats)
}

func (h *DashboardHandler) GetSalesReport(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	type SalesReport struct {
		TotalSales    float64 `json:"total_sales"`
		TotalOrders   int64   `json:"total_orders"`
		AverageOrder  float64 `json:"average_order"`
		TopCategories []struct {
			CategoryName string  `json:"category_name"`
			TotalSales   float64 `json:"total_sales"`
			TotalOrders  int64   `json:"total_orders"`
		} `json:"top_categories"`
	}

	var report SalesReport

	// Total sales and orders
	h.db.Model(&models.Transaction{}).
		Where("status = ? AND DATE(created_at) BETWEEN ? AND ?", "paid", startDate, endDate).
		Select("COALESCE(SUM(total), 0) as total_sales, COUNT(*) as total_orders").
		Scan(&report)

	// Average order value
	if report.TotalOrders > 0 {
		report.AverageOrder = report.TotalSales / float64(report.TotalOrders)
	}

	// Top categories
	h.db.Raw(`
		SELECT 
			categories.name as category_name,
			COALESCE(SUM(transaction_items.total_price), 0) as total_sales,
			COUNT(DISTINCT transactions.id) as total_orders
		FROM transaction_items
		JOIN menu_items ON transaction_items.menu_item_id = menu_items.id
		JOIN categories ON menu_items.category_id = categories.id
		JOIN transactions ON transaction_items.transaction_id = transactions.id
		WHERE transactions.status = 'paid' AND DATE(transactions.created_at) BETWEEN ? AND ?
		GROUP BY categories.id, categories.name
		ORDER BY total_sales DESC
		LIMIT 5
	`, startDate, endDate).Scan(&report.TopCategories)

	c.JSON(http.StatusOK, report)
}

func (h *DashboardHandler) GetProfitAnalysis(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	type ProfitAnalysis struct {
		GrossProfit  float64 `json:"gross_profit"`
		NetProfit    float64 `json:"net_profit"`
		ProfitMargin float64 `json:"profit_margin"`
		COGS         float64 `json:"cogs"`
		Revenue      float64 `json:"revenue"`
		Expenses     float64 `json:"expenses"`
		AddOnRevenue float64 `json:"addon_revenue"`
		AddOnCOGS    float64 `json:"addon_cogs"`
	}

	var analysis ProfitAnalysis

	// Revenue from paid transactions
	h.db.Model(&models.Transaction{}).
		Where("status = ? AND DATE(created_at) BETWEEN ? AND ?", "paid", startDate, endDate).
		Select("COALESCE(SUM(total), 0)").
		Scan(&analysis.Revenue)

	// COGS calculation for menu items
	h.db.Raw(`
		SELECT COALESCE(SUM(menu_items.cogs * transaction_items.quantity), 0) as cogs
		FROM transaction_items
		JOIN menu_items ON transaction_items.menu_item_id = menu_items.id
		JOIN transactions ON transaction_items.transaction_id = transactions.id
		WHERE transactions.status = 'paid' AND DATE(transactions.created_at) BETWEEN ? AND ?
	`, startDate, endDate).Scan(&analysis)

	// Add-on revenue and COGS
	h.db.Raw(`
		SELECT 
			COALESCE(SUM(transaction_item_add_ons.total_price), 0) as addon_revenue,
			COALESCE(SUM(add_ons.cogs * transaction_item_add_ons.quantity), 0) as addon_cogs
		FROM transaction_item_add_ons
		JOIN add_ons ON transaction_item_add_ons.add_on_id = add_ons.id
		JOIN transaction_items ON transaction_item_add_ons.transaction_item_id = transaction_items.id
		JOIN transactions ON transaction_items.transaction_id = transactions.id
		WHERE transactions.status = 'paid' AND DATE(transactions.created_at) BETWEEN ? AND ?
	`, startDate, endDate).Scan(&analysis)

	// Operational expenses
	h.db.Model(&models.Expense{}).
		Where("DATE(date) BETWEEN ? AND ?", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&analysis.Expenses)

	// Calculate total COGS including add-ons
	totalCOGS := analysis.COGS + analysis.AddOnCOGS

	// Calculate profits
	analysis.GrossProfit = analysis.Revenue - totalCOGS
	analysis.NetProfit = analysis.GrossProfit - analysis.Expenses

	// Profit margin
	if analysis.Revenue > 0 {
		analysis.ProfitMargin = (analysis.NetProfit / analysis.Revenue) * 100
	}

	c.JSON(http.StatusOK, analysis)
}
