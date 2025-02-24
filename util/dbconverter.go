package util

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math/rand"
	"pet-henry-3/data"
	"pet-henry-3/models"
)

// ** duplicate code because i dont want these exported ** //

// getRandomZeroNumber returns a random number padded with 0s
func getRandomZeroNumber() string {
	n := rand.Intn(1_000)
	return fmt.Sprintf("%04d", n)
}

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// generateSyncCode generates a random 5 digit 'syncCode' used for account recovery/syncing
func generateSyncCode() string {
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// ** end duplicate code ** //

// ConvertLegacy converts the SQLite DB users from the previous iteration of Pet Henry into the new MySQL database.
// It is only meant to be run once, but I'm keeping it here because its cool and i might need it again some day
func ConvertLegacy() {
	log.Println("Connecting to legacy database...")
	sqlite, err := sql.Open("sqlite3", "data.db")

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	defer sqlite.Close()

	log.Println("Connected to legacy database")
	log.Println("Retrieving legacy users...")
	rows, err := sqlite.Query("SELECT * FROM users")

	legacyUsers := []models.User{}

	for rows.Next() {
		legacyUser := models.User{}

		err := rows.Scan(&legacyUser.UserID, &legacyUser.PetCount, &legacyUser.DisplayName, &legacyUser.SyncCode)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}

		if legacyUser.SyncCode == "NEEDCODEPLS" {
			legacyUser.SyncCode = generateSyncCode()
		}

		fmt.Println(legacyUser)
		if legacyUser.PetCount > 0 {
			legacyUsers = append(legacyUsers, legacyUser)
		}
	}

	for _, user := range legacyUsers {

		_, err := data.DB.Exec("INSERT INTO users (user_id, display_name, sync_code, pet_count) VALUES (?, ?, ?, ?)", user.UserID, user.DisplayName, user.SyncCode, user.PetCount)

		if err != nil {
			log.Fatalf("Error inserting user: %v", err)
		}

		log.Println("Inserted user: ", user.DisplayName, "{", user.UserID, "} into new DB")
	}

	log.Println("Done converting legacy users...")
}
