package handlers

import (
	"btwarch/config"
	"btwarch/database"
	"btwarch/repositories"
	"btwarch/services"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	config         *config.Config
	githubService  *services.GitHubService
	jwtService     *services.JWTService
	userRepository *repositories.UserRepository
}

func NewAuthHandler(config *config.Config) *AuthHandler {
	githubService := services.NewGitHubService(
		config.GitHubClientID,
		config.GitHubClientSecret,
		config.GitHubRedirectURL,
	)
	jwtService := services.NewJWTService(config.JWTSecret)
	userRepository := repositories.NewUserRepository()

	return &AuthHandler{
		config:         config,
		githubService:  githubService,
		jwtService:     jwtService,
		userRepository: userRepository,
	}
}

func (h *AuthHandler) InitiateGitHubAuth(c *fiber.Ctx) error {
	state := generateRandomState()

	authURL := h.githubService.GetAuthURL(state)

	return c.Redirect(authURL)
}

func (h *AuthHandler) GitHubCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Authorization code is required",
		})
	}

	token, err := h.githubService.ExchangeCode(code)
	if err != nil {
		log.Printf("Error exchanging code for token: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to authenticate with GitHub",
		})
	}

	githubUser, err := h.githubService.GetUserInfo(token)
	if err != nil {
		log.Printf("Error getting user info: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user information",
		})
	}

	existingUser, err := h.userRepository.GetUserByGitHubID(githubUser.ID)
	if err != nil {
		log.Printf("Error checking existing user: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	var user *database.User
	if existingUser == nil {
		user, err = h.userRepository.CreateUser(
			githubUser.ID,
			githubUser.Login,
			githubUser.Email,
			githubUser.AvatarURL,
			token.AccessToken,
		)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}
	} else {
		err = h.userRepository.UpdateUserTokens(
			existingUser.ID.String(),
			token.AccessToken,
		)
		if err != nil {
			log.Printf("Error updating user tokens: %v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update user",
			})
		}
		user = existingUser
	}

	jwtToken, err := h.jwtService.GenerateToken(user.ID.String(), user.Username)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate authentication token",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Authentication successful",
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"avatar":   user.AvatarURL,
		},
		"token": jwtToken,
	})
}

func generateRandomState() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
