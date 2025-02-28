package data

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"pet-henry-3/models"
)

var ctx = context.Background()
var RDB *redis.Client

// InitCache initializes the redis connection
func InitCache() {
	log.Println("Initializing cache...")

	RDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to redis: %v", err)
		panic("Error connecting to redis")
	}

	loadFromDatabase()

	log.Println("Cache initialized")
}

func IncrementUserPet(u string) {
	res := RDB.ZIncrBy(ctx, "user_pets", 1, u)
	log.Println(res)
	printTopTen()
}

func printTopTen() {
	res, err := RDB.ZRevRangeWithScores(ctx, "user_pets", 0, 9).Result()
	fmt.Println(res, err)
}

// loadFromDatabase loads the user data from db into cache
func loadFromDatabase() {
	log.Println("Loading data into redis...")

	rows, err := DB.Query("SELECT user_id, display_name, sync_code, pet_count FROM users")
	if err != nil {
		log.Fatal("Error loading users into cache: ", err)
	}
	defer rows.Close()

	log.Println("Fetching users from database...")
	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.DisplayName, &user.SyncCode, &user.PetCount); err != nil {
			log.Fatal("Error scanning user rows: ", err)
		}
		users = append(users, user)
		_, rErr := RDB.ZAdd(ctx, "user_pets", redis.Z{Score: float64(user.PetCount), Member: user.UserID}).Result()

		if rErr != nil {
			log.Fatal("Error adding user to cache: ", rErr)
		}
	}

	res2, err := RDB.ZRevRangeWithScores(ctx, "user_pets", 0, 9).Result()
	log.Println(res2)

	log.Println("Data loaded into cache")
}
