package config

import (
	"os"
	"strings"
)

type Config struct {
	GitHubClientID     string
	GitHubClientSecret string
	GitHubRedirectURL  string

	CloudFlareZoneId   string
	CloudFlareApiToken string

	JWTSecret string

	DatabaseURL string

	Port string

	CookieDomain   string
	CookieSecure   bool
	CookieSameSite string

	ParentDomain string

	CORSOrigins []string
}

func getEnvArray(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	values := strings.Split(value, ",")
	for i := range values {
		values[i] = strings.TrimSpace(values[i])
	}

	return values
}

func LoadConfig() *Config {
	return &Config{
		GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		GitHubRedirectURL:  getEnv("GITHUB_REDIRECT_URL", "http://localhost:8080/auth/github/callback"),

		CloudFlareZoneId:   getEnv("CLOUDFLARE_ZONE_ID", ""),
		CloudFlareApiToken: getEnv("CLOUDFLARE_API_TOKEN", ""),

		JWTSecret: getEnv("JWT_SECRET", ""),

		DatabaseURL: getEnv("DATABASE_URL", "postgres://btwarch:btwarch@localhost:5432/btwarch?sslmode=disable"),

		Port: getEnv("PORT", "8080"),

		CookieDomain:   getEnv("COOKIE_DOMAIN", ""),
		CookieSecure:   getEnvBool("COOKIE_SECURE", false),
		CookieSameSite: getEnv("COOKIE_SAME_SITE", "lax"),

		ParentDomain: getEnv("PARENT_DOMAIN", "btwarch.me"),

		CORSOrigins: getEnvArray("CORS_ORIGINS", []string{}),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if value == "true" {
			return true
		}
		if value == "false" {
			return false
		}
	}
	return defaultValue
}
