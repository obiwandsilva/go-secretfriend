package friendsrepository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/obiwandsilva/go-secretfriend/application/config"
	"github.com/obiwandsilva/go-secretfriend/domain/entities"
)

type FriendsRepository struct {
	EnvConfig config.EnvironmentConfig
}

func New(envConfig config.EnvironmentConfig) *FriendsRepository {
	return &FriendsRepository{EnvConfig: envConfig}
}

func (fr *FriendsRepository) SaveFriendsList(friendsList []entities.Friend) error {
	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("error when starting aws session: %w", err)
	}

	uploader := s3manager.NewUploader(sess)
	jsonEncoded, err := json.Marshal(map[string][]entities.Friend{
		"friendsList": friendsList,
	})
	reader := bytes.NewReader(jsonEncoded)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(fr.EnvConfig.BucketName),
		Key: aws.String("friendsList.json"),
		Body: reader,
	})
	if err != nil {
		return fmt.Errorf("error when uploading file to S3: %w", err)
	}

	return nil
}

func (fr *FriendsRepository) GetFriendsList() ([]entities.Friend, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("error when starting aws session: %w", err)
	}

	buffer := make([]byte, 2048)
	bufferWriter := aws.NewWriteAtBuffer(buffer)

	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(bufferWriter,
		&s3.GetObjectInput{
			Bucket: aws.String(fr.EnvConfig.BucketName),
			Key: aws.String("friendsList.json"),
		})
	if err != nil {
		return nil, fmt.Errorf("error when downloading file from S3: %w", err)
	}

	friendsList := make(map[string][]entities.Friend)
	jsonReader := bytes.NewReader(buffer[0:numBytes])

	err = json.NewDecoder(jsonReader).Decode(&friendsList)
	if err != nil {
		return nil, fmt.Errorf("error when decoding file from S3: %w", err)
	}

	return friendsList["friendsList"], nil
}

func (fr *FriendsRepository) SaveDrawIndividualResult(phoneNumber string, picked entities.Friend) error {
	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("error when starting aws session: %w", err)
	}

	uploader := s3manager.NewUploader(sess)
	jsonEncoded, err := json.Marshal(picked)
	reader := bytes.NewReader(jsonEncoded)

	fileName := fmt.Sprintf("results/%s.json", phoneNumber)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(fr.EnvConfig.BucketName),
		Key: aws.String(fileName),
		Body: reader,
	})
	if err != nil {
		return fmt.Errorf("error when uploading file to S3: %w", err)
	}

	return nil
}

func (fr *FriendsRepository) SaveSecretKey(phoneNumber, secretKey string) error {
	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("error when starting aws session: %w", err)
	}
	
	s3svc := s3.New(sess)
	sourceKey := fmt.Sprintf("results/%s.json", phoneNumber)
	destKey := fmt.Sprintf("results/%s_%s.json", phoneNumber, secretKey)

	err = fr.copyS3Object(s3svc, sourceKey, destKey)
	if err != nil {
		return err
	}

	return nil
}

func (fr *FriendsRepository) GetSecretFriend(phoneNumber, secretKey string) (secretFriend entities.Friend, err error) {
	sess, err := session.NewSession()
	if err != nil {
		return secretFriend, fmt.Errorf("error when starting aws session: %w", err)
	}

	buffer := make([]byte, 516)
	bufferWriter := aws.NewWriteAtBuffer(buffer)
	downloader := s3manager.NewDownloader(sess)

	numBytes, err := downloader.Download(bufferWriter,
		&s3.GetObjectInput{
			Bucket: aws.String(fr.EnvConfig.BucketName),
			Key: aws.String(fmt.Sprintf("results/%s_%s.json", phoneNumber, secretKey)),
		})
	if err != nil {
		return secretFriend, fmt.Errorf("error when downloading file from S3: %w", err)
	}

	jsonReader := bytes.NewReader(buffer[0:numBytes])
	err = json.NewDecoder(jsonReader).Decode(&secretFriend)
	if err != nil {
		return secretFriend, fmt.Errorf("error when decoding file from S3: %w", err)
	}

	return secretFriend, nil
}

func (fr *FriendsRepository) copyS3Object(s3svc *s3.S3, sourceKey, destKey string) error {
	_, err := s3svc.CopyObject(&s3.CopyObjectInput{
		CopySource: aws.String(fmt.Sprintf("%s/%s", fr.EnvConfig.BucketName, sourceKey)),
		Bucket: aws.String(fr.EnvConfig.BucketName),
		Key: aws.String(destKey),
	})
	if err != nil {
		return fmt.Errorf("error when requesting object copy: %w", err)
	}

	// Wait to see if the item got copied
	err = s3svc.WaitUntilObjectExists(&s3.HeadObjectInput{
		Bucket: aws.String(fr.EnvConfig.BucketName),
		Key: aws.String(destKey),
	})
	if err != nil {
		return fmt.Errorf("error occurred while waiting for object %s to be copied to %s: %w",
			sourceKey,
			destKey,
			err,
		)
	}

	return nil
}

func (fr *FriendsRepository) deleteS3Object(s3svc *s3.S3, objectKey string) error {
	// Delete the object
	_, err := s3svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(fr.EnvConfig.BucketName),
		Key: aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("unable to delete object %s from bucket %q: %w", objectKey, fr.EnvConfig.BucketName, err)
	}

	err = s3svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(fr.EnvConfig.BucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf(
			"error occurred while waiting for object %s to be deleted from bucket %s: %w",
			objectKey,
			fr.EnvConfig.BucketName,
			err,
		)
	}

	return nil
}
