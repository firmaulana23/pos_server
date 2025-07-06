package models

import (
	"testing"
	"time"
)

func TestUserModel(t *testing.T) {
	user := User{
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "cashier",
		IsActive: true,
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username to be 'testuser', got %s", user.Username)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email to be 'test@example.com', got %s", user.Email)
	}

	if user.Role != "cashier" {
		t.Errorf("Expected role to be 'cashier', got %s", user.Role)
	}

	if !user.IsActive {
		t.Error("Expected user to be active")
	}
}

func TestMenuItemMarginCalculation(t *testing.T) {
	item := MenuItem{
		Name:  "Test Coffee",
		Price: 25000,
		COGS:  10000,
	}

	// Calculate margin (Price - COGS) / Price * 100
	expectedMargin := ((item.Price - item.COGS) / item.Price) * 100
	item.Margin = expectedMargin

	if item.Margin != 60.0 {
		t.Errorf("Expected margin to be 60.0, got %f", item.Margin)
	}
}

func TestAddOnMarginCalculation(t *testing.T) {
	addOn := AddOn{
		Name:  "Extra Shot",
		Price: 8000,
		COGS:  4000,
	}

	// Calculate margin
	expectedMargin := ((addOn.Price - addOn.COGS) / addOn.Price) * 100
	addOn.Margin = expectedMargin

	if addOn.Margin != 50.0 {
		t.Errorf("Expected margin to be 50.0, got %f", addOn.Margin)
	}
}

func TestTransactionModel(t *testing.T) {
	transaction := Transaction{
		TransactionNo: "TXN-001",
		Status:        "pending",
		SubTotal:      25000,
		Tax:           2500,
		Discount:      0,
		Total:         27500,
	}

	if transaction.TransactionNo != "TXN-001" {
		t.Errorf("Expected transaction number to be 'TXN-001', got %s", transaction.TransactionNo)
	}

	if transaction.Status != "pending" {
		t.Errorf("Expected status to be 'pending', got %s", transaction.Status)
	}

	expectedTotal := transaction.SubTotal + transaction.Tax - transaction.Discount
	if transaction.Total != expectedTotal {
		t.Errorf("Expected total to be %f, got %f", expectedTotal, transaction.Total)
	}
}

func TestExpenseModel(t *testing.T) {
	expense := Expense{
		Type:        "raw_material",
		Category:    "Coffee Beans",
		Description: "Premium coffee beans",
		Amount:      500000,
		Date:        time.Now(),
	}

	if expense.Type != "raw_material" {
		t.Errorf("Expected type to be 'raw_material', got %s", expense.Type)
	}

	if expense.Amount != 500000 {
		t.Errorf("Expected amount to be 500000, got %f", expense.Amount)
	}
}
