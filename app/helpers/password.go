package helpers

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var minPasswordLength = 8

func HashApiKeySecret(secret string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(secret), 6)
	return string(b), err
}

// ComparePassword checks if password and password hash match
func ComparePassword(passwordHash string, password string) error {
	passwordHashBytes := []byte(passwordHash)
	passwordBytes := []byte(password)
	return bcrypt.CompareHashAndPassword(passwordHashBytes, passwordBytes)
}

// GeneratePasswordHash generates the password hash from password
func GeneratePasswordHash(password string) (string, error) {
	saltedBytes := []byte(password)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, 12)
	if err != nil {
		return "", err
	}

	passwordHash := string(hashedBytes[:])
	return passwordHash, nil
}

// ValidatePassword checks if password has enough characters and contains only valid characters
func ValidatePassword(password string) error {

	if len(password) < minPasswordLength {
		return fmt.Errorf("Failed to validate password of length %d", len(password))
	}

	if regexp.MustCompile("[[:^graph:]]").FindStringIndex(password) != nil {
		return fmt.Errorf("Failed to validate password with invalid characters")
	}

	return nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// // ValidatePasswordResetToken checks whether a user has a token, and if they do, whether it's still valid
// func ValidatePasswordResetToken(user *entities.User) error {

// 	if user.ResetPasswordTokenExpiresAt.Valid && user.ResetPasswordTokenExpiresAt.Time.After(time.Now()) {
// 		return NewErrorWithCode(
// 			errors.New("valid reset password token exists"),
// 			ErrorCodePasswordResetTokenExists,
// 			"User has an existing valid password reset token.",
// 		)
// 	}

// 	return nil

// }
