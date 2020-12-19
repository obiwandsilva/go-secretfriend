package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/obiwandsilva/go-secretfriend/application/config"
	"github.com/obiwandsilva/go-secretfriend/domain/entities"
	"github.com/obiwandsilva/go-secretfriend/resources/repositories/friendsrepository"
)

type Event struct {
}

func handleRequest(_ context.Context, _ Event) (map[string][]entities.Friend, error) {
	envConfig := config.EnvironmentConfig{BucketName: "secret-friend"}
	friendsRepository := friendsrepository.New(envConfig)
	friendsList, err := friendsRepository.GetFriendsList()
	if err != nil {
		return nil, err
	}

	return map[string][]entities.Friend{"friendsList": friendsList}, nil
}

func main() {
	lambda.Start(handleRequest)
}

