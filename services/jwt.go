package services

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	secretKey      []byte
	cookieDomain   string
	cookieSecure   bool
	cookieSameSite string
}

type Claims struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	jwt.RegisteredClaims
}

func NewAuthService(secretKey string, cookieDomain string, cookieSecure bool, cookieSameSite string) *AuthService {
	return &AuthService{
		secretKey:      []byte(secretKey),
		cookieDomain:   cookieDomain,
		cookieSecure:   cookieSecure,
		cookieSameSite: cookieSameSite,
	}
}

func (a *AuthService) GenerateToken(userID string, username string, avatarURL string) (string, error) {
	claims := &Claims{
		UserID:    userID,
		Username:  username,
		AvatarURL: avatarURL,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secretKey)
}

func (a *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (a *AuthService) SetAuthCookie(c *fiber.Ctx, userID string, username string, avatarURL string) error {
	token, err := a.GenerateToken(userID, username, avatarURL)
	if err != nil {
		return err
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HTTPOnly = true
	cookie.Secure = a.cookieSecure
	cookie.SameSite = a.cookieSameSite
	cookie.Path = "/"

	if a.cookieDomain != "" {
		cookie.Domain = a.cookieDomain
	}

	c.Cookie(cookie)
	return nil
}

func (a *AuthService) ClearAuthCookie(c *fiber.Ctx) {
	cookie := new(fiber.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1 * time.Hour) // Expire immediately
	cookie.HTTPOnly = true
	cookie.Secure = a.cookieSecure
	cookie.SameSite = a.cookieSameSite
	cookie.Path = "/"

	if a.cookieDomain != "" {
		cookie.Domain = a.cookieDomain
	}

	c.Cookie(cookie)
}
