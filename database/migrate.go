package database

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Migration struct {
	Version   string
	Name      string
	SQL       string
	AppliedAt *sql.NullTime
}

func RunMigrations(db *sql.DB, migrationsDir string) error {
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	appliedMigrations, err := GetAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %v", err)
	}

	migrationFiles, err := ReadMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migration files: %v", err)
	}

	sort.Strings(migrationFiles)

	for _, file := range migrationFiles {
		version := ExtractVersion(file)
		if _, exists := appliedMigrations[version]; exists {
			log.Printf("Migration %s already applied, skipping", version)
			continue
		}

		if err := applyMigration(db, file, version); err != nil {
			return fmt.Errorf("failed to apply migration %s: %v", version, err)
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			version VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := db.Exec(query)
	return err
}

func GetAppliedMigrations(db *sql.DB) (map[string]Migration, error) {

	var tableExists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'migrations'
		);
	`).Scan(&tableExists)

	if err != nil {
		return nil, fmt.Errorf("failed to check if migrations table exists: %v", err)
	}

	if !tableExists {
		return make(map[string]Migration), nil
	}

	query := `SELECT version, name, applied_at FROM migrations ORDER BY version`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	migrations := make(map[string]Migration)
	for rows.Next() {
		var migration Migration
		err := rows.Scan(&migration.Version, &migration.Name, &migration.AppliedAt)
		if err != nil {
			return nil, err
		}
		migrations[migration.Version] = migration
	}

	return migrations, nil
}

func ReadMigrationFiles(migrationsDir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func ExtractVersion(filePath string) string {
	filename := filepath.Base(filePath)
	parts := strings.Split(filename, "_")
	if len(parts) >= 2 {
		return parts[0]
	}
	return filename
}

func applyMigration(db *sql.DB, filePath, version string) error {

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %v", err)
	}

	filename := filepath.Base(filePath)
	name := strings.TrimSuffix(filename, ".sql")

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute migration SQL: %v", err)
	}

	_, err = tx.Exec(
		"INSERT INTO migrations (version, name) VALUES ($1, $2)",
		version, name,
	)
	if err != nil {
		return fmt.Errorf("failed to record migration: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %v", err)
	}

	log.Printf("Applied migration %s: %s", version, name)
	return nil
}
