package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims dengan interface{} untuk fleksibilitas tipe data
type CustomClaims struct {
	Data map[string]interface{} `json:"data"`
	jwt.RegisteredClaims
}

var mySigningKey = []byte(os.Getenv("TOKEN_SALT"))

// ValidateToken untuk mengembalikan map[string]interface{}
func ValidateToken(myToken string) (bool, map[string]interface{}) {
	token, err := jwt.ParseWithClaims(myToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if err != nil {
		return false, nil
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return false, nil
	}

	return token.Valid, claims.Data
}

// ClaimToken untuk membuat token dengan data berbagai tipe
func ClaimToken(data map[string]interface{}, expirationHours int) (string, error) {
	if expirationHours <= 0 {
		expirationHours = 5
	}

	claims := CustomClaims{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expirationHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

// Helper function untuk mengkonversi tipe data
func GetString(claims map[string]interface{}, key string) string {
	if val, ok := claims[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func GetInt(claims map[string]interface{}, key string) int {
	if val, ok := claims[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case float32:
			return int(v)
		}
	}
	return 0
}

func GetBool(claims map[string]interface{}, key string) bool {
	if val, ok := claims[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

func GetFloat64(claims map[string]interface{}, key string) float64 {
	if val, ok := claims[key]; ok {
		if f, ok := val.(float64); ok {
			return f
		}
	}
	return 0.0
}
