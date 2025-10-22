package main

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "../data/inventory.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Generate hashes
	hash1, _ := bcrypt.GenerateFromPassword([]byte("qwertyasd"), bcrypt.DefaultCost)
	hash2, _ := bcrypt.GenerateFromPassword([]byte("luxperfume"), bcrypt.DefaultCost)

	// Insert users
	users := []struct {
		username, name, email, hash, role string
	}{
		{"Lux_perfume", "Lux Perfume Administrator", "lux@example.com", string(hash1), "admin"},
		{"manager", "Store Manager", "manager@example.com", string(hash2), "manager"},
	}

	for _, u := range users {
		_, err := db.Exec(`
			INSERT INTO Users (Username, Name, Email, PasswordHash, Role, IsActive)
			VALUES (?, ?, ?, ?, ?, 1)
		`, u.username, u.name, u.email, u.hash, u.role)

		if err != nil {
			log.Printf("Error creating %s: %v", u.username, err)
		} else {
			fmt.Printf("✓ Created user: %s (Role: %s)\n", u.username, u.role)
		}
	}
}
