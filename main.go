package main

import (
	"github.com/joho/godotenv"
	"log"
	"pet-henry-3/data"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file")
	}

	log.Println("Initializing PetHenry")

	data.InitDatabase()
	data.InitCache()
}
