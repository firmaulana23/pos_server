package database

import (
	"fmt"
	"log"
	"pos-system/internal/config"
	"pos-system/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.MenuItem{},
		&models.AddOn{},
		&models.Transaction{},
		&models.TransactionItem{},
		&models.TransactionItemAddOn{},
		&models.Expense{},
		&models.PaymentMethod{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Seed default data
	if err := seedDefaultData(db); err != nil {
		log.Printf("Warning: failed to seed default data: %v", err)
	}

	return &Database{DB: db}, nil
}

func seedDefaultData(db *gorm.DB) error {
	// Seed default payment methods
	paymentMethods := []models.PaymentMethod{
		{Name: "Cash", Code: "cash", IsActive: true},
		{Name: "Credit Card", Code: "card", IsActive: true},
		{Name: "Digital Wallet", Code: "digital_wallet", IsActive: true},
	}

	for _, pm := range paymentMethods {
		var existing models.PaymentMethod
		if err := db.Where("code = ?", pm.Code).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&pm).Error; err != nil {
					return fmt.Errorf("failed to create payment method %s: %w", pm.Name, err)
				}
			}
		}
	}

	// Seed default categories
	categories := []models.Category{
		// {Name: "Coffee", Description: "All types of coffee"},
		// {Name: "Tea", Description: "All types of tea"},
		// {Name: "Snacks", Description: "Light snacks and pastries"},
		// {Name: "Beverages", Description: "Non-coffee beverages"},
	}

	for _, cat := range categories {
		var existing models.Category
		if err := db.Where("name = ?", cat.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&cat).Error; err != nil {
					return fmt.Errorf("failed to create category %s: %w", cat.Name, err)
				}
			}
		}
	}

	// Seed default add-ons
	addOns := []models.AddOn{
		// {Name: "Extra Shot", Description: "Additional espresso shot", Price: 5000, COGS: 2000, IsAvailable: true},
		// {Name: "Whipped Cream", Description: "Fresh whipped cream", Price: 3000, COGS: 1500, IsAvailable: true},
		// {Name: "Extra Milk", Description: "Additional milk", Price: 2000, COGS: 1000, IsAvailable: true},
		// {Name: "Vanilla Syrup", Description: "Vanilla flavored syrup", Price: 3000, COGS: 1200, IsAvailable: true},
		// {Name: "Caramel Syrup", Description: "Caramel flavored syrup", Price: 3000, COGS: 1200, IsAvailable: true},
		// {Name: "Decaf", Description: "Decaffeinated option", Price: 0, COGS: 0, IsAvailable: true},
	}

	for _, addon := range addOns {
		var existing models.AddOn
		if err := db.Where("name = ?", addon.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&addon).Error; err != nil {
					return fmt.Errorf("failed to create add-on %s: %w", addon.Name, err)
				}
			}
		}
	}

	return nil
}
