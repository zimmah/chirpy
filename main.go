package main

import (
	"fmt"
	"os"

	"github.com/zimmah/chirpy/internal/database"
	"github.com/zimmah/chirpy/internal/router"
)

func main() {
	_, err := database.NewDB("./database.json")
	if err != nil {
		fmt.Println("Error creating database")
		os.Exit(1)
	}
	router.Router()
}