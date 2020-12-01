package main

import (
	"github.com/obiwandsilva/go-secretfriend/file"
	"log"
	"os"
)

func main() {
	log.Println("Starting raffle...")

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	log.Println(path)

	file.ReadFile("friends")
}
