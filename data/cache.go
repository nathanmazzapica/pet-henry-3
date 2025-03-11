package data

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"pet-henry-3/models"
	"strconv"
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
	// This here increments the score in the leaderboard cache
	res := RDB.ZIncrBy(ctx, "user_pets", 1, u)
	res2, err := RDB.HIncrBy(ctx, getHashKeyFromUID(u), "PetCount", 1).Result()
	log.Println(res)
	log.Println(res2)
	if err != nil {
		log.Fatalf("Error incrementing pet count: %v", err)
	}
	printTopTen()
}

func printTopTen() {
	res, err := RDB.ZRevRangeWithScores(ctx, "user_pets", 0, 9).Result()
	fmt.Println(res, err)
}

// loadFromDatabase loads the user data from db into leaderboard cache.
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
		log.Println("== Populating Leaderboard Set ==")
		_, rErr := RDB.ZAdd(ctx, "user_pets", redis.Z{Score: float64(user.PetCount), Member: user.UserID}).Result()

		if rErr != nil {
			log.Fatal("Error adding user to cache: ", rErr)
		}
	}

	res2, err := RDB.ZRevRangeWithScores(ctx, "user_pets", 0, 9).Result()
	log.Println(res2)

	log.Println("Data loaded into cache")
}

func InsertIntoHash(user *models.User) {
	hashFields := []string{
		"DisplayName", user.DisplayName,
		"PetCount", strconv.Itoa(user.PetCount),
		"SyncCode", user.SyncCode,
	}
	log.Println(hashFields)
	res, err := RDB.HSet(ctx, getHashKeyFromUID(user.UserID), hashFields).Result()

	if err != nil {
		log.Printf("== ERROR SETTING HASH FIELDS FOR USER %s ==\n", user.DisplayName)
		log.Println("Error setting hash fields: ", err)
		log.Println("==============================================================")
	}

	log.Println("Success:", res)
}

func GetFromHash(uid string) (*models.User, error) {
	// hacky quick thing to get key from uid
	key := getHashKeyFromUID(uid)
	exists, err := RDB.Exists(ctx, key).Result()

	if err != nil {
		log.Println("Error checking for user:", err)
	}

	if exists == 0 {
		log.Println("User does not exist", key)
		return nil, nil
	}

	log.Println("User found:", key)

	userData, err := RDB.HGetAll(ctx, key).Result()
	if err != nil {
		log.Println("Error getting user from database:", err)
	}

	fmt.Println(userData)

	user := &models.User{}

	user.UserID = uid
	user.DisplayName = userData["DisplayName"]
	log.Println("PET COUNT:", userData["PetCount"])
	user.PetCount, err = strconv.Atoi(userData["PetCount"])
	log.Println(user.PetCount)
	user.SyncCode = userData["SyncCode"]

	return user, err
}

type playerScore struct {
	user     string `redis:"user"`
	petCount int64
}

func GetTopPlayersWithScores() []playerScore {
	var players []playerScore
	if res, err := RDB.ZRevRangeWithScores(ctx, "user_pets", 0, 9).Result(); err == nil {
		for _, entry := range res {
			fmt.Println(entry.Member, entry.Score)
			players = append(players, playerScore{user: entry.Member.(string), petCount: int64(entry.Score)})
		}
	}

	return players
}

func getHashKeyFromUID(uid string) string {
	return fmt.Sprintf("users:%s", uid)
}
