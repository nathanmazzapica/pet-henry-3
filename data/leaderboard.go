package data

import (
	"pet-henry-3/models"
)

var (
	topTen = make(map[int]models.User)
)

func InitLeaderboard() {
	_ = GetTopPlayersWithScores()
}
