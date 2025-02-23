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

	log.Println("Initializing PetHenry")

	data.InitDatabase()
	data.InitPetCounter()

	data.InitCache()

	server.InitRoutes()

	fmt.Println("===============\nPET DAISY STARTED\n===============\n\n")
	server.StartHTTP()

}
