package main

import (
	"btwarch/config"
	"btwarch/database"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	if len(os.Args) > 1 && os.Args[1] == "status" {
		showStatus()
		return
	}

	runMigrations()
}

func runMigrations() {
	cfg := config.LoadConfig()

	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "./database/migrations"
	}

	log.Printf("Running migrations from: %s", migrationsDir)

	if err := database.RunMigrations(database.DB, migrationsDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
}

func showStatus() {
	cfg := config.LoadConfig()

	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "./database/migrations"
	}

	fmt.Printf("Migration Status for: %s\n", migrationsDir)
	fmt.Println("=" + strings.Repeat("=", len(migrationsDir)+20))

	appliedMigrations, err := database.GetAppliedMigrations(database.DB)
	if err != nil {
		log.Fatalf("Failed to get applied migrations: %v", err)
	}

	migrationFiles, err := database.ReadMigrationFiles(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migration files: %v", err)
	}

	sort.Strings(migrationFiles)

	fmt.Printf("%-15s %-30s %-20s\n", "Version", "Name", "Status")
	fmt.Println(strings.Repeat("-", 70))

	for _, file := range migrationFiles {
		version := database.ExtractVersion(file)
		filename := filepath.Base(file)
		name := strings.TrimSuffix(filename, ".sql")

		status := "Pending"
		if _, exists := appliedMigrations[version]; exists {
			status = "Applied"
		}

		fmt.Printf("%-15s %-30s %-20s\n", version, name, status)
	}

	fmt.Println()
	fmt.Printf("Total migrations: %d\n", len(migrationFiles))
	fmt.Printf("Applied: %d\n", len(appliedMigrations))
	fmt.Printf("Pending: %d\n", len(migrationFiles)-len(appliedMigrations))
}
