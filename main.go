package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/zimmah/chirpy/internal/database"
	"github.com/zimmah/chirpy/internal/router"
)

func main() {
	godotenv.Load()

	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *debug {
		fmt.Println("Debug mode enabled. Deleting database.json...")

		err := os.Remove("database.json")
		if err != nil {
			fmt.Printf("Error deleting file: %v\n", err)
		} else {
			fmt.Println("database.json deleted succesfully.")
		}
	}
	_, err := database.NewDB("./database.json")
	if err != nil {
		fmt.Println("Error creating database")
		os.Exit(1)
	}
	router.Router()
}