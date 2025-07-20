package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net"
	"net/http"
	"strings"
	"templ_workout/internals/models"

	"golang.org/x/crypto/bcrypt"
)

var (
	userContextKey = "user"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("failed to generate token: %v", err)
	}
	return base64.RawStdEncoding.EncodeToString(bytes)
}

func GetClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}

	realIP := r.Header.Get("X-Real-Ip")
	if realIP != "" {
		return realIP
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // fallback, may include port
	}
	return ip
}

func GetUserFromContext(r *http.Request) *models.User {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}

func GetUserEmailFromContext(ctx context.Context) string {
	result := ""
	if user, ok := ctx.Value("user").(*models.User); ok && user != nil {
		result = user.Name
	}
	return result
}

func IsLoggedIn(ctx context.Context) bool {
	if user, ok := ctx.Value("user").(*models.User); ok && user != nil {
		return true
	}
	return false
}
