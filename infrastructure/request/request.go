package request

import "errors"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type PasswordMatcher interface {
	ArePasswordsMatch() error
}

type RegistrationRequest struct {
	Email           string `json:"email" validate:"required,email"`
	FirstName       string `json:"firstName" validate:"required"`
	LastName        string `json:"lastName" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
	Address         string `json:"address" validate:"required"`
	Group           string `json:"group" validate:"required"`
}

type UpdatingUserRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Address   string `json:"address" validate:"required"`
}

type UpdatePasswordRequest struct {
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
}

type VerificationRequest struct {
	VerificationId string `json:"verificationId" validate:"required"`
	UserId         string `json:"userId" validate:"required"`
	SecurityCode   int    `json:"securityCode" validate:"required,numeric,min=1000,max=9999"`
}

type SendVerificationCodeRequest struct {
	VerificationId string `json:"verificationId" validate:"required"`
	UserEmail      string `json:"userEmail" validate:"required,email"`
}

func (r RegistrationRequest) ArePasswordsMatch() error {
	return checkPasswordsMatch(r.Password, r.ConfirmPassword)
}

func (r UpdatePasswordRequest) ArePasswordsMatch() error {
	return checkPasswordsMatch(r.Password, r.ConfirmPassword)
}

func checkPasswordsMatch(password, confirmPassword string) error {
	if password != confirmPassword {
		return errors.New("passwords do not match")
	}
	return nil
}
