package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("error while loading env file")
	}

	// database := db.New(context.Background())
	// log.Println("connected to database")
}
