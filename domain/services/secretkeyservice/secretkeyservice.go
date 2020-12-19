package secretkeyservice

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/obiwandsilva/go-secretfriend/application/config"
	"github.com/obiwandsilva/go-secretfriend/resources/repositories/friendsrepository"
	"strings"
)

type SecretKeyService struct {
	EnvConfig config.EnvironmentConfig
	FriendsRepository *friendsrepository.FriendsRepository
}

func New(
	envConfig config.EnvironmentConfig,
	friendsRepository *friendsrepository.FriendsRepository,
) *SecretKeyService {
	return &SecretKeyService{EnvConfig: envConfig, FriendsRepository: friendsRepository}
}

func (sks *SecretKeyService) GenerateSecretKey(phoneNumber string) (string, error) {
	created, err := sks.secretKeyAlreadyCreatedForPhoneNumber(phoneNumber)
	if err != nil {
		return "", err
	}
	if created {
		return "", NewSecretKeyAlreadyCreatedError(phoneNumber)
	}

	secretKey := strings.ToUpper(uuid.New().String()[:4])

	err = sks.FriendsRepository.SaveSecretKey(phoneNumber, secretKey)
	if err != nil {
		return "", fmt.Errorf("error while saving secretkey for %s: %w", phoneNumber, err)
	}

	return secretKey, nil
}

func (sks *SecretKeyService) secretKeyAlreadyCreatedForPhoneNumber(phoneNumber string) (bool, error) {
	sess, err := session.NewSession()
	if err != nil {
		return false, fmt.Errorf("error when starting aws session: %w", err)
	}

	s3svc := s3.New(sess)
	objectKeyPrefix := fmt.Sprintf("results/%s_", phoneNumber)
	output, err := s3svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:              aws.String(sks.EnvConfig.BucketName),
		Prefix:              aws.String(objectKeyPrefix),
	})
	if err != nil {
		return false, fmt.Errorf("error when listing objects: %w", err)
	}

	if len(output.Contents) == 0 {
		return false, nil
	}

	return true, nil
}