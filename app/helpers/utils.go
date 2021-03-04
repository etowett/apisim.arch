package helpers

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/google/uuid"
)

// SliceContains checks whether a string is in a given struct
func SliceContains(givenSlice []string, a string) bool {
	for _, b := range givenSlice {
		if b == a {
			return true
		}
	}
	return false
}

// GetUUID returns a uniqueue UUID string
func GetUUID() string {
	return uuid.New().String()
}

// GetMD5Hash returns a uniqueue md5 string
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
