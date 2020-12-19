package sms

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
)

type TwilioInfo struct {
	sid string
	authToken string
	from string
}

func SendMessage(message, phoneNumber string) error {
	return sendThroughSNS(message, phoneNumber)
}

func sendThroughSNS(message, phoneNumber string) error {
	sess, err := session.NewSession()
	if err != nil {
		log.Println("could not start aws session")
		return fmt.Errorf("error when startin aws session: %w", err)
	}

	snsSvc := sns.New(sess)

	params := &sns.PublishInput{
		Message:                aws.String(message),
		PhoneNumber:            aws.String(phoneNumber),
	}
	_, err = snsSvc.Publish(params)
	if err != nil {
		return fmt.Errorf("could not send SMS to friend '%s' with number '%s': %w", err)
	}

	return nil
}
