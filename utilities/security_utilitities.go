package utilities

import (
	"crypto/rand"
	"fmt"
)

// GenerateRandomToken generated 32 bit long
func GenerateRandomToken() string {
	bytes := make([]byte, 16)
	_, errorInfo := rand.Read(bytes)
	var token string
	if errorInfo == nil {
		token = fmt.Sprintf("%x", bytes)
	}
	return token
}
