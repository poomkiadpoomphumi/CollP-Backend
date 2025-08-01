package middleware
import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"github.com/golang-jwt/jwt/v5"
	"crypto/rsa"
	"log"
)

// Store public key in package variable or inject via function
var publicKey *rsa.PublicKey
// Call this once from your app init with your loaded RSA public key
func SetPublicKey(key *rsa.PublicKey) {
	publicKey = key
}
func extractTokenFromHeader(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}
	return parts[1], nil
}
func AuthMiddleware(next http.Handler) http.Handler {
	if publicKey == nil {
		log.Fatal("public key is not set in AuthMiddleware")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}
		tokenString, err := extractTokenFromHeader(authHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is RS256
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
