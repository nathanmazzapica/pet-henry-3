package data

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"os"
	"pet-henry-3/game"
	"pet-henry-3/models"
)

var DB *sql.DB

// InitDatabase initializes the database connection and loads the
// data into the redis cache
func InitDatabase() {
	log.Println("Initializing database...")

	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASS"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_HOST"),
		DBName: os.Getenv("DB_NAME"),
	}

	var err error
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalf("Error initializing database: ")
		panic(err)
	}

	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatalf("Error connecting to the database: ")
		panic(pingErr)
	}

	log.Println("Database initialized")
}

func InitPetCounter() {
	result := DB.QueryRow("SELECT SUM(pet_count) FROM users")
	result.Scan(&game.Counter)
	log.Println("GetCurrentCount: ", game.Counter)
}

func CreateNewUser() *models.User {
	userID := uuid.New().String()
	displayName := getRandomDisplayName()

	newUser := &models.User{userID, displayName, generateSyncCode(), 0}

	_, err := DB.Exec("INSERT INTO users (user_id, display_name, sync_code, pet_count) VALUES (?, ?, ?, ?)", newUser.UserID, newUser.DisplayName, newUser.SyncCode, newUser.PetCount)

	if err != nil {
		log.Fatalf("Error creating new user: %v", err)
	}

	return newUser
}

func GetUserFromDB(userID string) (*models.User, error) {
	user := &models.User{}

	result := DB.QueryRow("SELECT user_id, display_name, sync_code, pet_count FROM users WHERE user_id = ?", userID)

	if err := result.Scan(&user.UserID, &user.DisplayName, &user.SyncCode, &user.PetCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id {%v} not found", userID)
		}
		return nil, fmt.Errorf("error scanning rows: %v", err)
	}
	return user, nil
}

// This and all other User specific functions will be moved to data/user.go
// This includes:
// 1. GetUserFromDB
// 2. CreateNewUser
// 3. getRandomDisplayName()
// 4. generateSyncCode
// 5. SaveToDB
func getRandomDisplayName() string {
	adjectives := []string{"big", "long", "small", "golden", "yellow", "black",
		"red", "short", "cunning", "silly", "radical", "sluggish",
		"speedy", "humorous", "shy", "scared", "brave", "intelligent", "stupid",
		"orange", "medium", "austere", "gaudy", "ugly", "beautiful", "sexy",
		"intellectual", "philosophical", "charged", "empty", "full",
		"serious", "vengeful", "malignant", "generous", "complacent",
		"ambitious", "lazy", "dull", "sharp", "splendid", "sexy", "cute",
		"loving", "hateful", "spiteful", "rude", "polite", "dasterdly"}

	nouns := []string{"Dog", "Watermelon", "Crusader", "Lancer", "Envisage", "Frog",
		"Beetle", "Cellphone", "Python", "Lizard", "Butterfly", "Dragon",
		"Automobile", "Cow", "Henry", "Levi", "Array", "Buzzer", "Balloon", "Book",
		"Calendar", "Burrito", "Corgi", "Pencil", "Pen", "Marker", "Bookshelf",
		"Sharpener", "Can", "Lightbulb", "Flower", "Daisy", "Eraser", "Battery",
		"Butter", "Cantaloupe", "Fridge", "Computer", "Programmer", "Kitty"}

	fmt.Printf("%d\n", len(adjectives)*len(nouns)*1_000)

	adjI := rand.Intn(len(adjectives))
	nounI := rand.Intn(len(nouns))

	return adjectives[adjI] + nouns[nounI] + getRandomZeroNumber()

}

// getRandomZeroNumber returns a random number padded with 0s
func getRandomZeroNumber() string {
	n := rand.Intn(1_000)
	return fmt.Sprintf("%04d", n)
}

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// generateSyncCode generates a random 5 digit 'syncCode' used for account recovery/syncing
func generateSyncCode() string {
	code := make([]byte, 5)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// autosave saves the DB on a set interval
func autosave() {

}
