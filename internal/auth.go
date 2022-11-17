package internal

import (
	"net/http"
	"strings"

	"github.com/amakmurr/dans-multi-pro-test/pkg/jwt"
)

type JWTAuthentication struct {
	jwtClient jwt.Client
}

func NewJWTAuthentication(jwtClient jwt.Client) JWTAuthentication {
	return JWTAuthentication{
		jwtClient: jwtClient,
	}
}

func (j *JWTAuthentication) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string
		bearer := r.Header.Get("Authorization")
		if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
			tokenString = bearer[7:]
		}
		if !j.jwtClient.VerifyToken(tokenString) {
			handleError(w, r, ErrUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
