package data

import (
	"log"
	"pet-henry-3/models"
)

var (
	topTen = make(map[int]models.User)
)

type LeaderboardRow struct {
	Position    int    `json:"position"` // maybe unnecessary?
	PetCount    int64  `json:"pet_count"`
	DisplayName string `json:"display_name"`
}

func InitLeaderboard() {
	log.Println("Initializing leaderboard")
	top := GetTopPlayersWithScores()

	for i, user := range top {
		log.Println("Index:", i, "User:", user.user)

		newRow := LeaderboardRow{Position: i + 1, PetCount: 69, DisplayName: "WIP"}

		log.Println(newRow)
	}
}
