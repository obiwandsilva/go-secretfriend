package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/obiwandsilva/go-secretfriend/application/config"
	"github.com/obiwandsilva/go-secretfriend/domain/services/templateservice"
	"github.com/obiwandsilva/go-secretfriend/resources/repositories/friendsrepository"
	"log"
	"net/http"
)

func handleRequest(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("handling get secret friend request")

	res := events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body: fmt.Sprintf("invalid phone number or secret key"),
	}

	phoneNumber, okPhoneNumber := req.PathParameters["phoneNumber"]
	secretKey, okSecretKey := req.PathParameters["secretKey"]

	log.Printf("phoneNumber: %s secretKey: %s\n", phoneNumber, secretKey)

	if okPhoneNumber && okSecretKey {
		friendsRepository := friendsrepository.New(config.EnvironmentConfig{BucketName: "secret-friend"})
		secretFriend, err := friendsRepository.GetSecretFriend(phoneNumber, secretKey)
		if err != nil {
			log.Printf("error when getting secret friend for number %s: %v\n", phoneNumber, secretFriend)
			res.Body = err.Error()
			res.StatusCode = http.StatusForbidden
			return res, err
		}
		log.Printf("secret friend for phone number %s: %v\n", phoneNumber, secretFriend)

		content, err := templateservice.RenderSecretFriend(phoneNumber, secretFriend)
		if err != nil {
			log.Printf("error when rendering content : %v\n", secretFriend)
			res.Body = err.Error()
			res.StatusCode = http.StatusInternalServerError
			return res, err
		}

		log.Println("RENDERED CONTENT: ", content)

		res.Headers = map[string]string{"Content-Type": "text/html"}
		res.StatusCode = http.StatusOK
		res.Body = content
	}

	return res, nil
}

func main() {
	lambda.Start(handleRequest)
}
