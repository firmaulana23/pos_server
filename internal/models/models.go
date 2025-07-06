package models

import (
	"time"
	"gorm.io/gorm"
)

// User represents users in the system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;default:''"`
	FullName  string         `json:"full_name" gorm:"default:''"`
	Password  string         `json:"-" gorm:"not null"`
	Role      string         `json:"role" gorm:"not null;default:'cashier'"` // admin, manager, cashier
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Category represents menu categories
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	MenuItems   []MenuItem     `json:"menu_items,omitempty"`
}

// MenuItem represents menu items
type MenuItem struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CategoryID  uint           `json:"category_id"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Price       float64        `json:"price" gorm:"not null"`
	COGS        float64        `json:"cogs" gorm:"not null"` // Cost of Goods Sold (HPP)
	Margin      float64        `json:"margin" gorm:"-"`      // Calculated field
	IsAvailable bool           `json:"is_available" gorm:"default:true"`
	ImageURL    string         `json:"image_url"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Category    Category       `json:"category,omitempty"`
}

// AddOn represents available add-ons for menu items
type AddOn struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Price       float64        `json:"price" gorm:"not null"`
	COGS        float64        `json:"cogs" gorm:"not null"`
	Margin      float64        `json:"margin" gorm:"-"` // Calculated field
	IsAvailable bool           `json:"is_available" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Transaction represents sales transactions
type Transaction struct {
	ID            uint                `json:"id" gorm:"primaryKey"`
	TransactionNo string              `json:"transaction_no" gorm:"uniqueIndex;not null"`
	UserID        uint                `json:"user_id"`
	Status        string              `json:"status" gorm:"not null;default:'pending'"` // pending, paid
	PaymentMethod string              `json:"payment_method"`                           // cash, card, digital_wallet
	SubTotal      float64             `json:"sub_total" gorm:"not null"`
	Tax           float64             `json:"tax" gorm:"default:0"`
	Discount      float64             `json:"discount" gorm:"default:0"`
	Total         float64             `json:"total" gorm:"not null"`
	PaidAt        *time.Time          `json:"paid_at"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	DeletedAt     gorm.DeletedAt      `json:"-" gorm:"index"`
	User          User                `json:"user,omitempty"`
	Items         []TransactionItem   `json:"items,omitempty"`
}

// TransactionItem represents items in a transaction
type TransactionItem struct {
	ID            uint                      `json:"id" gorm:"primaryKey"`
	TransactionID uint                      `json:"transaction_id"`
	MenuItemID    uint                      `json:"menu_item_id"`
	Quantity      int                       `json:"quantity" gorm:"not null"`
	UnitPrice     float64                   `json:"unit_price" gorm:"not null"`
	TotalPrice    float64                   `json:"total_price" gorm:"not null"`
	CreatedAt     time.Time                 `json:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at"`
	MenuItem      MenuItem                  `json:"menu_item,omitempty"`
	Transaction   Transaction               `json:"transaction,omitempty"`
	AddOns        []TransactionItemAddOn    `json:"add_ons,omitempty"`
}

// TransactionItemAddOn represents add-ons for transaction items
type TransactionItemAddOn struct {
	ID                uint            `json:"id" gorm:"primaryKey"`
	TransactionItemID uint            `json:"transaction_item_id"`
	AddOnID           uint            `json:"add_on_id"`
	Quantity          int             `json:"quantity" gorm:"not null;default:1"`
	UnitPrice         float64         `json:"unit_price" gorm:"not null"`
	TotalPrice        float64         `json:"total_price" gorm:"not null"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	AddOn             AddOn           `json:"add_on,omitempty"`
	TransactionItem   TransactionItem `json:"transaction_item,omitempty"`
}

// Expense represents business expenses
type Expense struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Type        string         `json:"type" gorm:"not null"` // raw_material, operational
	Category    string         `json:"category" gorm:"not null"`
	Description string         `json:"description" gorm:"not null"`
	Amount      float64        `json:"amount" gorm:"not null"`
	Date        time.Time      `json:"date" gorm:"not null"`
	UserID      uint           `json:"user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	User        User           `json:"user,omitempty"`
}

// PaymentMethod represents available payment methods
type PaymentMethod struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Code      string    `json:"code" gorm:"uniqueIndex;not null"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
