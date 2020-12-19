package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/obiwandsilva/go-secretfriend/application/config"
	"github.com/obiwandsilva/go-secretfriend/domain/entities"
	"github.com/obiwandsilva/go-secretfriend/resources/repositories/friendsrepository"
)

type Event struct {
	FriendsList []entities.Friend `json:"friendsList"`
}

func handleRequest(_ context.Context, event Event) (string, error) {
	envConfig := config.EnvironmentConfig{BucketName: "secret-friend"}
	friendsRepository := friendsrepository.New(envConfig)
	err := friendsRepository.SaveFriendsList(event.FriendsList)
	if err != nil {
		return err.Error(), err
	}

	return "OK", nil
}

func main() {
	lambda.Start(handleRequest)
}
