package secretkeyservice

import "fmt"

type SecretKeyAlreadyCreatedError struct {
	Message string
	PhoneNumber string
}

func NewSecretKeyAlreadyCreatedError(phoneNumber string) *SecretKeyAlreadyCreatedError {
	return &SecretKeyAlreadyCreatedError{
		Message: fmt.Sprintf("A secret key was already generated for number %s", phoneNumber),
		PhoneNumber: phoneNumber,
	}
}

func (e *SecretKeyAlreadyCreatedError) Error() string {
	return e.Message
}
