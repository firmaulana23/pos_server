package auth

import (
	"testing"
)

func TestGenerateAndValidateToken(t *testing.T) {
	secretKey := "test-secret-key"
	expiryHours := 1
	jwtService := NewJWTService(secretKey, expiryHours)
	
	userID := uint(1)
	username := "testuser"
	role := "cashier"

	// Generate token
	token, err := jwtService.GenerateToken(userID, username, role)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Generated token should not be empty")
	}

	// Validate token
	claims, err := jwtService.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID to be %d, got %d", userID, claims.UserID)
	}

	if claims.Username != username {
		t.Errorf("Expected Username to be %s, got %s", username, claims.Username)
	}

	if claims.Role != role {
		t.Errorf("Expected Role to be %s, got %s", role, claims.Role)
	}
}

func TestValidateInvalidToken(t *testing.T) {
	secretKey := "test-secret-key"
	expiryHours := 1
	jwtService := NewJWTService(secretKey, expiryHours)
	
	invalidToken := "invalid.token.here"

	_, err := jwtService.ValidateToken(invalidToken)
	if err == nil {
		t.Error("Expected error when validating invalid token")
	}
}

func TestHashAndVerifyPassword(t *testing.T) {
	password := "testpassword123"

	// Hash password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hashedPassword == "" {
		t.Error("Hashed password should not be empty")
	}

	if hashedPassword == password {
		t.Error("Hashed password should be different from original password")
	}

	// Verify correct password
	if err := CheckPassword(hashedPassword, password); err != nil {
		t.Error("Password verification should succeed with correct password")
	}

	// Verify incorrect password
	if err := CheckPassword(hashedPassword, "wrongpassword"); err == nil {
		t.Error("Password verification should fail with incorrect password")
	}
}

func TestTokenClaims(t *testing.T) {
	claims := &Claims{
		UserID:   123,
		Username: "john_doe",
		Role:     "admin",
	}

	if claims.UserID != 123 {
		t.Errorf("Expected UserID to be 123, got %d", claims.UserID)
	}

	if claims.Username != "john_doe" {
		t.Errorf("Expected Username to be 'john_doe', got %s", claims.Username)
	}

	if claims.Role != "admin" {
		t.Errorf("Expected Role to be 'admin', got %s", claims.Role)
	}
}
