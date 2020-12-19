package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/obiwandsilva/go-secretfriend/application/config"
	"github.com/obiwandsilva/go-secretfriend/domain/entities"
	"github.com/obiwandsilva/go-secretfriend/domain/services/raffleservice"
	"github.com/obiwandsilva/go-secretfriend/resources/repositories/friendsrepository"
)

var friendsRepository = friendsrepository.New(config.EnvironmentConfig{BucketName: "secret-friend"})

func generateObjects(pairs entities.Pairs) error {
	for picker, picked := range pairs {
		err := friendsRepository.SaveDrawIndividualResult(picker.PhoneNumber, picked)
		if err != nil {
			return fmt.Errorf("error when saving individual result: %w", err)
		}
	}

	return nil
}

func handleRequest(_ context.Context) (string, error) {
	friendsList, err := friendsRepository.GetFriendsList()
	if err != nil {
		return "", fmt.Errorf("error when getting friends list: %w", err)
	}

	pool := raffleservice.GeneratePool(friendsList)
	pairs := entities.Pairs{}

	raffleservice.Draw(friendsList, pool, pairs)

	err = generateObjects(pairs)
	if err != nil {
		return "", err
	}

	return "OK", nil
}

func main() {
	lambda.Start(handleRequest)
}
