package data

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
	"pet-henry-3/game"
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

// autosave saves the DB on a set interval
func autosave() {

}
