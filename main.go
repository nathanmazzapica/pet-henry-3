package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"pet-henry-3/data"
	"pet-henry-3/server"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}

	log.Println("Initializing Daisy")

	data.InitDatabase()
	data.InitPetCounter()

	data.InitCache()
	data.InitLeaderboard()

	//util.ConvertLegacy()

	server.InitRoutes()
	server.InitBroadcasts()

	fmt.Println("\n\n========\n           _____  ______ _______     _____          _____  _______     __    _____ _______       _____ _______ ______ _____           \n  ______  |  __ \\|  ____|__   __|   |  __ \\   /\\   |_   _|/ ____\\ \\   / /   / ____|__   __|/\\   |  __ \\__   __|  ____|  __ \\   ______ \n |______| | |__) | |__     | |      | |  | | /  \\    | | | (___  \\ \\_/ /   | (___    | |  /  \\  | |__) | | |  | |__  | |  | | |______|\n  ______  |  ___/|  __|    | |      | |  | |/ /\\ \\   | |  \\___ \\  \\   /     \\___ \\   | | / /\\ \\ |  _  /  | |  |  __| | |  | |  ______ \n |______| | |    | |____   | |      | |__| / ____ \\ _| |_ ____) |  | |      ____) |  | |/ ____ \\| | \\ \\  | |  | |____| |__| | |______|\n          |_|    |______|  |_|      |_____/_/    \\_\\_____|_____/   |_|     |_____/   |_/_/    \\_\\_|  \\_\\ |_|  |______|_____/          \n                                                                                                                                     \n===============\n\n")
	server.StartHTTP()

}
