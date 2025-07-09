package main

import (
	"fmt"
	"log"
	"pos-system/internal/config"
	"pos-system/internal/database"
	"pos-system/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("🌱 Seeding database with test data...")

	// Create users
	if err := seedUsers(db); err != nil {
		log.Printf("Failed to seed users: %v", err)
	}

	// Create categories
	// if err := seedCategories(db); err != nil {
	// 	log.Printf("Failed to seed categories: %v", err)
	// }

	// Create menu items
	// if err := seedMenuItems(db); err != nil {
	// 	log.Printf("Failed to seed menu items: %v", err)
	// }

	// Create add-ons
	// if err := seedAddOns(db); err != nil {
	// 	log.Printf("Failed to seed add-ons: %v", err)
	// }

	// Create expenses
	// if err := seedExpenses(db); err != nil {
	// 	log.Printf("Failed to seed expenses: %v", err)
	// }

	fmt.Println("✅ Database seeding completed!")
	fmt.Println("\n📝 Test Accounts:")
	fmt.Println("   Admin: admin@pos.com / admin123")
	fmt.Println("   Manager: manager@pos.com / manager123")
	fmt.Println("   Cashier: cashier@pos.com / cashier123")
}

func seedUsers(db *database.Database) error {
	users := []models.User{
		{
			Username: "admin",
			Email:    "admin@pos.com",
			Role:     "admin",
			IsActive: true,
		},
		{
			Username: "manager",
			Email:    "manager@pos.com",
			Role:     "manager",
			IsActive: true,
		},
		{
			Username: "cashier",
			Email:    "cashier@pos.com",
			Role:     "cashier",
			IsActive: true,
		},
	}

	passwords := []string{"admin123", "manager123", "cashier123"}

	for i, user := range users {
		// Check if user already exists
		var existingUser models.User
		if err := db.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
			fmt.Printf("⚠️  User %s already exists, skipping\n", user.Email)
			continue
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwords[i]), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)

		if err := db.DB.Create(&user).Error; err != nil {
			return err
		}
		fmt.Printf("✅ Created user: %s (%s)\n", user.Email, user.Role)
	}

	return nil
}

func seedCategories(db *database.Database) error {
	categories := []models.Category{
		{Name: "Coffee", Description: "Hot and cold coffee beverages"},
		{Name: "Tea", Description: "Various tea selections"},
		{Name: "Pastries", Description: "Fresh baked goods"},
		{Name: "Sandwiches", Description: "Light meals and sandwiches"},
		{Name: "Beverages", Description: "Non-coffee drinks"},
	}

	for _, category := range categories {
		var existing models.Category
		if err := db.DB.Where("name = ?", category.Name).First(&existing).Error; err == nil {
			fmt.Printf("⚠️  Category %s already exists, skipping\n", category.Name)
			continue
		}

		if err := db.DB.Create(&category).Error; err != nil {
			return err
		}
		fmt.Printf("✅ Created category: %s\n", category.Name)
	}

	return nil
}

func seedMenuItems(db *database.Database) error {
	// Get category IDs
	var coffeeCategory, teaCategory, pastryCategory models.Category
	db.DB.Where("name = ?", "Coffee").First(&coffeeCategory)
	db.DB.Where("name = ?", "Tea").First(&teaCategory)
	db.DB.Where("name = ?", "Pastries").First(&pastryCategory)

	menuItems := []models.MenuItem{
		// Coffee items
		{CategoryID: coffeeCategory.ID, Name: "Espresso", Description: "Rich and strong coffee shot", Price: 15000, COGS: 8000, IsAvailable: true},
		{CategoryID: coffeeCategory.ID, Name: "Americano", Description: "Espresso with hot water", Price: 18000, COGS: 10000, IsAvailable: true},
		{CategoryID: coffeeCategory.ID, Name: "Cappuccino", Description: "Espresso with steamed milk and foam", Price: 25000, COGS: 12000, IsAvailable: true},
		{CategoryID: coffeeCategory.ID, Name: "Latte", Description: "Espresso with steamed milk", Price: 28000, COGS: 14000, IsAvailable: true},
		{CategoryID: coffeeCategory.ID, Name: "Iced Coffee", Description: "Cold brewed coffee with ice", Price: 22000, COGS: 11000, IsAvailable: true},

		// Tea items
		{CategoryID: teaCategory.ID, Name: "Green Tea", Description: "Premium green tea", Price: 15000, COGS: 7000, IsAvailable: true},
		{CategoryID: teaCategory.ID, Name: "Earl Grey", Description: "Classic black tea with bergamot", Price: 16000, COGS: 8000, IsAvailable: true},
		{CategoryID: teaCategory.ID, Name: "Chamomile", Description: "Relaxing herbal tea", Price: 18000, COGS: 9000, IsAvailable: true},

		// Pastries
		{CategoryID: pastryCategory.ID, Name: "Croissant", Description: "Buttery and flaky pastry", Price: 12000, COGS: 6000, IsAvailable: true},
		{CategoryID: pastryCategory.ID, Name: "Muffin", Description: "Blueberry muffin", Price: 15000, COGS: 7500, IsAvailable: true},
		{CategoryID: pastryCategory.ID, Name: "Danish", Description: "Sweet Danish pastry", Price: 18000, COGS: 9000, IsAvailable: true},
	}

	for _, item := range menuItems {
		var existing models.MenuItem
		if err := db.DB.Where("name = ? AND category_id = ?", item.Name, item.CategoryID).First(&existing).Error; err == nil {
			fmt.Printf("⚠️  Menu item %s already exists, skipping\n", item.Name)
			continue
		}

		if err := db.DB.Create(&item).Error; err != nil {
			return err
		}
		fmt.Printf("✅ Created menu item: %s (Rp %.0f)\n", item.Name, item.Price)
	}

	return nil
}

func seedAddOns(db *database.Database) error {
	addOns := []models.AddOn{
		{Name: "Extra Shot", Description: "Additional espresso shot", Price: 8000, COGS: 4000, IsAvailable: true},
		{Name: "Decaf", Description: "Decaffeinated option", Price: 0, COGS: 0, IsAvailable: true},
		{Name: "Soy Milk", Description: "Replace with soy milk", Price: 5000, COGS: 3000, IsAvailable: true},
		{Name: "Oat Milk", Description: "Replace with oat milk", Price: 7000, COGS: 4000, IsAvailable: true},
		{Name: "Extra Hot", Description: "Served extra hot", Price: 0, COGS: 0, IsAvailable: true},
		{Name: "Extra Foam", Description: "Additional milk foam", Price: 3000, COGS: 1500, IsAvailable: true},
		{Name: "Vanilla Syrup", Description: "Sweet vanilla flavoring", Price: 5000, COGS: 2500, IsAvailable: true},
		{Name: "Caramel Syrup", Description: "Sweet caramel flavoring", Price: 5000, COGS: 2500, IsAvailable: true},
	}

	for _, addOn := range addOns {
		var existing models.AddOn
		if err := db.DB.Where("name = ?", addOn.Name).First(&existing).Error; err == nil {
			fmt.Printf("⚠️  Add-on %s already exists, skipping\n", addOn.Name)
			continue
		}

		if err := db.DB.Create(&addOn).Error; err != nil {
			return err
		}
		fmt.Printf("✅ Created add-on: %s (Rp %.0f)\n", addOn.Name, addOn.Price)
	}

	return nil
}

func seedExpenses(db *database.Database) error {
	// Get admin user ID
	var admin models.User
	db.DB.Where("role = ?", "admin").First(&admin)

	expenses := []models.Expense{
		{Type: "raw_material", Category: "Coffee Beans", Description: "Premium Arabica coffee beans - 5kg", Amount: 500000, Date: time.Now().AddDate(0, 0, -5), UserID: admin.ID},
		{Type: "raw_material", Category: "Milk", Description: "Fresh milk - 20L", Amount: 150000, Date: time.Now().AddDate(0, 0, -4), UserID: admin.ID},
		{Type: "raw_material", Category: "Sugar", Description: "White sugar - 10kg", Amount: 75000, Date: time.Now().AddDate(0, 0, -3), UserID: admin.ID},
		{Type: "raw_material", Category: "Cups & Lids", Description: "Paper cups and lids - 500pcs", Amount: 200000, Date: time.Now().AddDate(0, 0, -2), UserID: admin.ID},
		{Type: "operational", Category: "Utilities", Description: "Electricity bill", Amount: 800000, Date: time.Now().AddDate(0, 0, -1), UserID: admin.ID},
		{Type: "operational", Category: "Rent", Description: "Shop rent", Amount: 3000000, Date: time.Now().AddDate(0, 0, -1), UserID: admin.ID},
		{Type: "operational", Category: "Staff", Description: "Staff salary", Amount: 2500000, Date: time.Now().AddDate(0, 0, -1), UserID: admin.ID},
	}

	for _, expense := range expenses {
		if err := db.DB.Create(&expense).Error; err != nil {
			return err
		}
		fmt.Printf("✅ Created expense: %s (Rp %.0f)\n", expense.Description, expense.Amount)
	}

	return nil
}
