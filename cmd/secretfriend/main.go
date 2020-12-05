package main

import (
	"fmt"
	"github.com/obiwandsilva/go-secretfriend/domain"
	"github.com/obiwandsilva/go-secretfriend/file"
	"github.com/obiwandsilva/go-secretfriend/sms"
	"log"
)

func main() {
	log.Println("Starting raffle...")

	filePath := "friends"

	log.Printf("reading friends list from file %s\n", filePath)
	friendsList, err := file.ReadFile(filePath)
	if err != nil {
		log.Panicf("could not read the file: %v", err)
	}

	log.Println("pool created")
	pool := generatePool(friendsList)
	pairs := domain.Pairs{}

	log.Println("drawing participants")
	pairs = domain.Draw(friendsList, pool, pairs)

	log.Println("sending SMS")
	err = sendMessages(pairs)
	if err != nil {
		log.Fatal(err)
	}
}

func generatePool(friendsList []domain.Friend) (pool []domain.Friend) {
	for _, f := range friendsList {
		pool = append(pool, f)
	}

	return pool
}

func printPairs(pairs domain.Pairs) {
	for picker, picked := range pairs {
		fmt.Printf("Picker: %s\nPicked: %s\n####################\n", picker.Name, picked.Name)
	}
}

func sendMessages(pairs domain.Pairs) error {
	message := "Roi, %s, ne? N são 4i20 mas ta na hra de saber seu ou sua sorteadx. Segue o nome. Guarde bem: %s"
	for picker, picked := range pairs {
		log.Printf("sending sms to %s %s", picker.Name, picker.PhoneNumber)
		err := sms.SendMessage(
			fmt.Sprintf(message, picker, picked),
			picker.PhoneNumber,
		)
		if err != nil {
			return fmt.Errorf(
				"error when sending SMS to %s with number %s: %w",
				picker.Name,
				picker.PhoneNumber,
				err,
			)
		}
	}

	return nil
}
