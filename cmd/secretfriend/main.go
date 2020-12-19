package main

import (
	"fmt"
	"github.com/obiwandsilva/go-secretfriend/domain/entities"
	"github.com/obiwandsilva/go-secretfriend/domain/services/raffleservice"
	"github.com/obiwandsilva/go-secretfriend/file"
	"github.com/obiwandsilva/go-secretfriend/resources/sms"
	"log"
	"os"
)

func main() {
	log.Println("Starting raffle...")

	filePath := os.Getenv("SECRET_FRIENDS_FILE")

	log.Printf("reading friends list from file %s\n", filePath)
	friendsList, err := file.ReadFile(filePath, true)
	if err != nil {
		log.Panicf("could not read the file: %v", err)
	}

	pool := raffleservice.GeneratePool(friendsList)
	log.Println("pool created")

	pairs := entities.Pairs{}

	log.Println("drawing participants")
	pairs = raffleservice.Draw(friendsList, pool, pairs)

	printPairs(pairs)
}

func printPairs(pairs entities.Pairs) {
	log.Println("printing pairs")
	for picker, picked := range pairs {
		message := fmt.Sprintf(
			"Roi. %s, né? Já passou de 16h20 mas tá na hora de saber seu ou sua sorteadx. Segue o nome. Guarde bem: %s",
			picker.Name,
			picked.Name,
		)

		fmt.Printf(
			"Picker: %s\nPicked: %s\nMessage: %s\n####################\n",
			picker.Name,
			picked.Name,
			message,
		)
	}
}

func sendMessages(pairs entities.Pairs) error {
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
