package helpers

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const (
	digits           = "0123456789"
	lowerCaseLetters = "abcdefghijklmnopqrstuvwxyz"
	upperCaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// func GenerateEmailActivationKey() string {
// 	return generateRandomString(digits+lowerCaseLetters+upperCaseLetters, 8)
// }

// func GeneratePhoneActivationKey() string {
// 	return generateRandomString(digits, 5)
// }

func GenerateApiKeyID() string {
	return generateRandomString(digits+lowerCaseLetters+upperCaseLetters, 20)
}

func GenerateApiKeySecret() string {
	return generateRandomString(digits+lowerCaseLetters+upperCaseLetters, 40)
}

// func GenerateSalt() string {
// 	return generateRandomString(digits+lowerCaseLetters+upperCaseLetters, 12)
// }

func GenerateUUID() string {
	return uuid.New().String()
}

func generateRandomString(charset string, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
