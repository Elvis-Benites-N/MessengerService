package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a hashed password from the given plain-text password.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err) // Add error context for better error handling
	}

	return string(hashedPassword), nil
}

// CheckPassword compares a plain-text password with a hashed password to validate if they match.
func CheckPassword(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("incorrect password: %w", err) // Add error context for better error handling
	}

	return nil
}
