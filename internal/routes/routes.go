package routes

import (
	"pos-system/internal/handlers"
	"pos-system/internal/middleware"
	"pos-system/pkg/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, jwtService *auth.JWTService) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, jwtService)
	menuHandler := handlers.NewMenuHandler(db)
	addOnHandler := handlers.NewAddOnHandler(db)
	transactionHandler := handlers.NewTransactionHandler(db)
	expenseHandler := handlers.NewExpenseHandler(db)
	dashboardHandler := handlers.NewDashboardHandler(db)

	// Apply CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Public routes
	api := router.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
		}

		// Public menu routes (for POS display)
		public := api.Group("/public")
		{
			public.GET("/menu/categories", menuHandler.GetCategories)
			public.GET("/menu/items", menuHandler.GetMenuItems)
			public.GET("/menu/items/:id", menuHandler.GetMenuItem)
			public.GET("/add-ons", addOnHandler.GetAddOns)
			public.GET("/add-ons/:id", addOnHandler.GetAddOn)
			public.GET("/payment-methods", transactionHandler.GetPaymentMethods)
		}
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtService))
	{
		// Profile routes
		profile := protected.Group("/profile")
		{
			profile.GET("", authHandler.GetProfile)
			profile.PUT("", authHandler.UpdateProfile)
		}

		// Menu management routes
		menu := protected.Group("/menu")
		{
			// Categories
			menu.GET("/categories", menuHandler.GetCategories)
			menu.POST("/categories", middleware.RequireRole("admin", "manager"), menuHandler.CreateCategory)
			menu.PUT("/categories/:id", middleware.RequireRole("admin", "manager"), menuHandler.UpdateCategory)
			menu.DELETE("/categories/:id", middleware.RequireRole("admin", "manager"), menuHandler.DeleteCategory)

			// Menu items
			menu.GET("/items", menuHandler.GetMenuItems)
			menu.GET("/items/:id", menuHandler.GetMenuItem)
			menu.POST("/items", middleware.RequireRole("admin", "manager"), menuHandler.CreateMenuItem)
			menu.PUT("/items/:id", middleware.RequireRole("admin", "manager"), menuHandler.UpdateMenuItem)
			menu.DELETE("/items/:id", middleware.RequireRole("admin", "manager"), menuHandler.DeleteMenuItem)
		}

		// Add-on management routes
		addOns := protected.Group("/add-ons")
		{
			addOns.GET("", addOnHandler.GetAddOns)
			addOns.GET("/:id", addOnHandler.GetAddOn)
			addOns.POST("", middleware.RequireRole("admin", "manager"), addOnHandler.CreateAddOn)
			addOns.PUT("/:id", middleware.RequireRole("admin", "manager"), addOnHandler.UpdateAddOn)
			addOns.DELETE("/:id", middleware.RequireRole("admin", "manager"), addOnHandler.DeleteAddOn)
		}

		// Transaction routes
		transactions := protected.Group("/transactions")
		{
			transactions.GET("", transactionHandler.GetTransactions)
			transactions.GET("/:id", transactionHandler.GetTransaction)
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.PUT("/:id/pay", transactionHandler.PayTransaction)
			transactions.DELETE("/:id", middleware.RequireRole("admin", "manager"), transactionHandler.DeleteTransaction)
		}

		// Payment methods
		protected.GET("/payment-methods", transactionHandler.GetPaymentMethods)

		// Expense routes
		expenses := protected.Group("/expenses")
		{
			expenses.GET("", expenseHandler.GetExpenses)
			expenses.GET("/:id", expenseHandler.GetExpense)
			expenses.POST("", middleware.RequireRole("admin", "manager"), expenseHandler.CreateExpense)
			expenses.PUT("/:id", middleware.RequireRole("admin", "manager"), expenseHandler.UpdateExpense)
			expenses.DELETE("/:id", middleware.RequireRole("admin", "manager"), expenseHandler.DeleteExpense)
			expenses.GET("/summary", middleware.RequireRole("admin", "manager"), expenseHandler.GetExpenseSummary)
		}

		// Dashboard routes
		dashboard := protected.Group("/dashboard")
		dashboard.Use(middleware.RequireRole("admin", "manager"))
		{
			dashboard.GET("/stats", dashboardHandler.GetDashboardStats)
			dashboard.GET("/sales-report", dashboardHandler.GetSalesReport)
			dashboard.GET("/profit-analysis", dashboardHandler.GetProfitAnalysis)
		}

		// User management routes (admin only)
		users := protected.Group("/users")
		users.Use(middleware.RequireRole("admin"))
		{
			users.GET("", authHandler.GetUsers)
			users.GET("/:id", authHandler.GetUser)
			users.PUT("/:id", authHandler.UpdateUser)
			users.PUT("/:id/role", authHandler.UpdateUserRole)
			users.DELETE("/:id", authHandler.DeleteUser)
		}
	}

	// Serve static files for admin dashboard
	router.Static("/static", "./web/static")
	router.LoadHTMLGlob("web/templates/*")

	// Admin dashboard routes
	admin := router.Group("/admin")
	{
		admin.GET("/", func(c *gin.Context) {
			c.HTML(200, "login.html", gin.H{
				"title": "POS System - Login",
			})
		})
		admin.GET("/dashboard", func(c *gin.Context) {
			c.HTML(200, "dashboard.html", gin.H{
				"title": "POS System - Dashboard",
			})
		})
		admin.GET("/menu", func(c *gin.Context) {
			c.HTML(200, "menu.html", gin.H{
				"title": "POS System - Menu Management",
			})
		})
		admin.GET("/add-ons", func(c *gin.Context) {
			c.HTML(200, "addons.html", gin.H{
				"title": "POS System - Add-ons Management",
			})
		})
		admin.GET("/transactions", func(c *gin.Context) {
			c.HTML(200, "transactions.html", gin.H{
				"title": "POS System - Transactions",
			})
		})
		admin.GET("/expenses", func(c *gin.Context) {
			c.HTML(200, "expenses.html", gin.H{
				"title": "POS System - Expenses",
			})
		})
		admin.GET("/pos", func(c *gin.Context) {
			c.HTML(200, "pos.html", gin.H{
				"title": "POS System - Point of Sale",
			})
		})
		admin.GET("/users", func(c *gin.Context) {
			c.HTML(200, "users.html", gin.H{
				"title": "POS System - User Management",
			})
		})
	}
}
