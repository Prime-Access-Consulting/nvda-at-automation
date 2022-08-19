package main

import (
	"AT/GolangServer/AT"
	"AT/GolangServer/server"
	"github.com/joho/godotenv"
	"log"
)

func loadEnv() {
	log.Println("Loading .env file;")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file!")
	}
}

func main() {
	loadEnv()
	new(server.Server).Start(AT.LoadATs())
	log.Println("Ready.")
}
