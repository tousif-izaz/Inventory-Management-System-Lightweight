package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open database
	db, err := sql.Open("sqlite3", "../../data/inventory.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Hash password "admin"
	password := "admin"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Insert admin user
	query := `
		INSERT OR REPLACE INTO Users (Username, Name, Email, PasswordHash, Role, IsActive, CreatedAt, UpdatedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	_, err = db.Exec(query,
		"admin",
		"Admin User",
		"admin@example.com",
		string(hashedPassword),
		"admin",
		1, // IsActive = true
		now,
		now,
	)

	if err != nil {
		log.Fatal("Failed to insert user:", err)
	}

	fmt.Println("✓ Admin user created successfully!")
	fmt.Println("  Username: admin")
	fmt.Println("  Password: admin")
	fmt.Println("\nPlease change the password after first login.")
}
