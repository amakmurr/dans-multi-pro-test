package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

type Client struct {
	secretKey string
}

func NewClient(secretKey string) Client {
	return Client{
		secretKey: secretKey,
	}
}

func (c *Client) GenerateToken(uc UserClaims) (string, error) {
	uc.StandardClaims.ExpiresAt = time.Now().Add(24 * time.Hour).UnixMilli()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenString, err := token.SignedString([]byte(c.secretKey))
	if err != nil {
		return "Signing Error", err
	}
	return tokenString, nil
}

func (c *Client) VerifyToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.secretKey), nil
	})
	if err != nil {
		return false
	}
	return token.Valid
}
