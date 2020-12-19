package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/obiwandsilva/go-secretfriend/application/config"
	"github.com/obiwandsilva/go-secretfriend/domain/services/secretkeyservice"
	"github.com/obiwandsilva/go-secretfriend/domain/services/templateservice"
	"github.com/obiwandsilva/go-secretfriend/resources/repositories/friendsrepository"
	"net/http"
)

func handleRequest(_ context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	res := events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body: "invalid phone number",
	}

	res.Headers = map[string]string{"Content-Type": "text/html"}

	if phoneNumber, ok := req.PathParameters["phoneNumber"]; ok {
		envConfig := config.EnvironmentConfig{BucketName: "secret-friend"}
		friendsRepository := friendsrepository.New(envConfig)
		secretKeyService := secretkeyservice.New(envConfig, friendsRepository)

		secretKey, err := secretKeyService.GenerateSecretKey(phoneNumber)
		if err != nil {
			res.Body = err.Error()

			var secretKeyAlreadyCreatedError *secretkeyservice.SecretKeyAlreadyCreatedError
			if errors.As(err, &secretKeyAlreadyCreatedError) {
				res.StatusCode = http.StatusCreated
				content, err := templateservice.RenderSecretKey(phoneNumber, "", true)
				if err != nil {
					res.Body = err.Error()
					return res, err
				}

				res.Body = content
				return res, nil
			}

			res.Body = "<h1>Erro ao gerar senha</h1>"
			res.StatusCode = http.StatusInternalServerError
			return res, err
		}

		content, err := templateservice.RenderSecretKey(phoneNumber, secretKey, false)
		if err != nil {
			res.Body = fmt.Sprintf("<h1>Erro ao gerar senha: %v</h1>", err)
			res.StatusCode = http.StatusInternalServerError
			return res, err
		}

		res.Body = content
		res.StatusCode = http.StatusCreated
		return res, err
	}

	return res, nil
}

func main() {
	lambda.Start(handleRequest)
}
