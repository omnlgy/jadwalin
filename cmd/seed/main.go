package main

import (
	"log"

	"github.com/gofrs/uuid/v5"
	"github.com/omnlgy/jadwalin/internal/config"
	"github.com/omnlgy/jadwalin/internal/db"
	"github.com/omnlgy/jadwalin/internal/domain"
	"github.com/omnlgy/jadwalin/internal/models"
)

func main() {
	// Load environment variables from .env file
	config := config.Load()

	db, err := db.NewPostgresDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database successfully.")

	// Auto-migrate tables (creates table if not exists, no-op if already exists)
	log.Println("Migrating database...")
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}
	log.Println("Database migrated.")

	// Check if users table already has data
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Printf("Users table already has %d records, skipping seed.\n", count)
		return
	}

	// Seed users
	log.Println("Seeding users...")
	users := []models.User{
		{
			PhoneNumber: "6281234567890",
			Email:       "john@example.com",
			Address:     "Jl. Sudirman No. 1",
			FullName:    "John Doe",
			Photo:       "http://example.com/john.jpg",
			Role:        string(domain.RoleAdmin),
		},
		{
			PhoneNumber: "6281122334455",
			Email:       "jane@example.com",
			Address:     "Jl. Thamrin No. 2",
			FullName:    "Jane Smith",
			Photo:       "http://example.com/jane.jpg",
			Role:        string(domain.RoleEmployee),
		},
		{
			PhoneNumber: "6287654321098",
			Email:       "peter@example.com",
			Address:     "Jl. Gatot Subroto No. 3",
			FullName:    "Peter Jones",
			Photo:       "http://example.com/peter.jpg",
			Role:        string(domain.RoleUser),
		},
	}

	for i := range users {
		// Manually generate UUIDv7 for seeding for consistency and to ensure it's set
		// even if BeforeCreate hook is somehow bypassed or if we want specific UUIDs
		newUUID, err := uuid.NewV7()
		if err != nil {
			log.Fatalf("Failed to generate UUIDv7: %v", err)
		}
		users[i].ID = newUUID

		if err := db.Create(&users[i]).Error; err != nil {
			log.Fatalf("Failed to create user %s: %v", users[i].FullName, err)
		}
		log.Printf("Created user: %s (ID: %s)\n", users[i].FullName, users[i].ID)
	}

	log.Println("Seeding completed successfully!")
}
