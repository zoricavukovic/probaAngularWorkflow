package application

import (
	"encoding/json"
	"github.com/mmmajder/zms-devops-auth-service/infrastructure/dto"
	"log"
	"net/http"
)

type UserService struct {
	HttpClient           *http.Client
	AuthService          *AuthService
	KeycloakService      *KeycloakService
	IdentityProviderHost string
}

func NewUserService(httpClient *http.Client, authService *AuthService, keycloakService *KeycloakService, identityProviderHost string) *UserService {

	return &UserService{
		HttpClient:           httpClient,
		AuthService:          authService,
		KeycloakService:      keycloakService,
		IdentityProviderHost: identityProviderHost,
	}
}

func (service *UserService) GetUser(authorizationHeader string) (*dto.UserDTO, error) {
	responseBody, err := service.KeycloakService.GetKeycloakUser(authorizationHeader)

	if err != nil {
		return nil, err
	}
	userDTO := &dto.UserDTO{}
	log.Print(responseBody)
	json.NewDecoder(responseBody).Decode(userDTO)
	return userDTO, nil
}

func (service *UserService) GetUserById(authorizationHeader, id string) (*dto.UserDTO, error) {
	log.Printf("GETTING USER ID: %s", id)
	responseBody, err := service.KeycloakService.GetKeycloakUserById(authorizationHeader, id)
	if err != nil {
		return nil, err
	}

	keycloakDTO := &dto.GettingKeycloakUserDTO{}
	json.NewDecoder(responseBody).Decode(keycloakDTO)
	log.Println(keycloakDTO)
	log.Println(keycloakDTO.Attributes)
	log.Println(keycloakDTO.Attributes.Address)

	return &dto.UserDTO{
		Id:         id,
		GivenName:  keycloakDTO.FirstName,
		FamilyName: keycloakDTO.LastName,
		Email:      keycloakDTO.Email,
		Address:    keycloakDTO.Attributes.Address[0],
	}, nil
}

func (service *UserService) UpdateUser(authorizationHeader, id, firstName, lastName, address string) error {
	requestBody := dto.NewUpdateKeycloakUserDTO(id, firstName, lastName, address)
	responseBody, err := service.KeycloakService.UpdateKeycloakUser(authorizationHeader, id, requestBody)

	if err != nil {
		return err
	}

	log.Println(responseBody)

	return nil
}

func (service *UserService) DeleteUser(authorizationHeader, id string) error {
	err := service.KeycloakService.DeleteKeycloakUser(authorizationHeader, id)
	if err != nil {
		return err
	}

	return nil
}

func (service *UserService) ResetPassword(authorizationHeader, id, password string) error {
	requestBody := dto.NewCredentialsDTO(password)

	err := service.KeycloakService.ResetPasswordOfKeycloakUser(authorizationHeader, id, requestBody)
	if err != nil {
		return err
	}

	return nil
}
