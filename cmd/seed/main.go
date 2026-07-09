package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/omnlgy/jadwalin/internal/config"
	"github.com/omnlgy/jadwalin/internal/db"
	"github.com/omnlgy/jadwalin/internal/domain"
	"github.com/omnlgy/jadwalin/internal/models"
)

func main() {
	cfg := config.Load()

	db, err := db.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to database successfully.")

	log.Println("Migrating database...")
	if err := db.AutoMigrate(&models.User{}, &models.Treatment{}, &models.Booking{}, &models.StaffSkill{}); err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}
	log.Println("Database migrated.")

	// Seed users
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		log.Printf("Users table already has %d records, skipping seed.\n", userCount)
	} else {
		log.Println("Seeding users...")
		users := []models.User{
			{
				PhoneNumber: "6281234567890",
				Email:       "john@example.com",
				Address:     "Jl. Sudirman No. 1",
				FullName:    "John Doe",
				Photo:       "http://example.com/john.jpg",
				Role:        string(domain.RoleAdmin),
				Verified:    true,
			},
			{
				PhoneNumber: "6281122334455",
				Email:       "jane@example.com",
				Address:     "Jl. Thamrin No. 2",
				FullName:    "Jane Smith",
				Photo:       "http://example.com/jane.jpg",
				Role:        string(domain.RoleStaff),
				Verified:    true,
			},
			{
				PhoneNumber: "6287654321098",
				Email:       "peter@example.com",
				Address:     "Jl. Gatot Subroto No. 3",
				FullName:    "Peter Jones",
				Photo:       "http://example.com/peter.jpg",
				Role:        string(domain.RoleUser),
				Verified:    true,
			},
		}

		for i := range users {
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
	}

	// Seed treatments
	var treatmentCount int64
	db.Model(&models.Treatment{}).Count(&treatmentCount)
	if treatmentCount > 0 {
		log.Printf("Treatments table already has %d records, skipping seed.\n", treatmentCount)
	} else {
		log.Println("Seeding treatments...")
		treatments := []models.Treatment{
			{
				Name:        "Haircut",
				Description: "Standard haircut for men and women.",
				Duration:    30,
				Price:       50.00,
			},
			{
				Name:        "Massage",
				Description: "Relaxing full body massage.",
				Duration:    60,
				Price:       100.00,
			},
			{
				Name:        "Manicure",
				Description: "Professional manicure and nail care.",
				Duration:    45,
				Price:       30.00,
			},
		}

		for i := range treatments {
			newUUID, err := uuid.NewV7()
			if err != nil {
				log.Fatalf("Failed to generate UUIDv7 for treatment: %v", err)
			}
			treatments[i].ID = newUUID

			if err := db.Create(&treatments[i]).Error; err != nil {
				log.Fatalf("Failed to create treatment %s: %v", treatments[i].Name, err)
			}
			log.Printf("Created treatment: %s (ID: %s)\n", treatments[i].Name, treatments[i].ID)
		}
	}

	// Seed bookings
	var bookingCount int64
	db.Model(&models.Booking{}).Count(&bookingCount)
	if bookingCount > 0 {
		log.Printf("Bookings table already has %d records, skipping seed.\n", bookingCount)
	} else {
		log.Println("Seeding bookings...")
		var users []models.User
		db.Find(&users)
		var treatments []models.Treatment
		db.Find(&treatments)

		if len(users) < 3 || len(treatments) < 2 {
			log.Println("Not enough users or treatments to seed bookings. Skipping booking seed.")
		} else {
			bookings := []models.Booking{
				{
					ClientID:    users[2].ID,      // Peter Jones
					StaffID:     users[1].ID,      // Jane Smith
					TreatmentID: treatments[0].ID, // Haircut
					StartTime:   time.Now().Add(24 * time.Hour),
					EndTime:     time.Now().Add(24 * time.Hour).Add(30 * time.Minute),
					Status:      "confirmed",
				},
				{
					ClientID:    users[2].ID,      // Peter Jones
					StaffID:     users[1].ID,      // Jane Smith
					TreatmentID: treatments[1].ID, // Massage
					StartTime:   time.Now().Add(48 * time.Hour),
					EndTime:     time.Now().Add(48 * time.Hour).Add(60 * time.Minute),
					Status:      "pending",
				},
			}

			for i := range bookings {
				if err := db.Create(&bookings[i]).Error; err != nil {
					log.Fatalf("Failed to create booking for client %s: %v", bookings[i].ClientID, err)
				}
				log.Printf("Created booking for client %s (ID: %s)\n", bookings[i].ClientID, bookings[i].ID)
			}
		}
	}

	// Seed staff skills
	var staffSkillCount int64
	db.Model(&models.StaffSkill{}).Count(&staffSkillCount)
	if staffSkillCount > 0 {
		log.Printf("StaffSkills table already has %d records, skipping seed.\n", staffSkillCount)
	} else {
		log.Println("Seeding staff skills...")
		var users []models.User
		db.Find(&users)
		var treatments []models.Treatment
		db.Find(&treatments)

		var staffUser models.User
		for _, u := range users {
			if u.Role == string(domain.RoleStaff) {
				staffUser = u
				break
			}
		}

		if staffUser.ID == uuid.Nil || len(treatments) == 0 {
			log.Println("No staff user or treatments found. Skipping staff skill seed.")
		} else {
			for _, t := range treatments {
				skill := models.StaffSkill{
					UserID:      staffUser.ID,
					TreatmentID: t.ID,
				}
				if err := db.Create(&skill).Error; err != nil {
					log.Printf("Failed to assign skill %s -> %s: %v\n", staffUser.FullName, t.Name, err)
				} else {
					log.Printf("Assigned skill: %s -> %s (ID: %s)\n", staffUser.FullName, t.Name, skill.ID)
				}
			}
		}
	}

	log.Println("Seeding completed successfully!")
}