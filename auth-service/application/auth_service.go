package application

import (
	"encoding/json"
	"errors"
	"github.com/mmmajder/zms-devops-auth-service/domain"
	"github.com/mmmajder/zms-devops-auth-service/infrastructure/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type AuthService struct {
	store           domain.VerificationStore
	HttpClient      *http.Client
	KeycloakService *KeycloakService
}

func NewAuthService(store domain.VerificationStore, httpClient *http.Client, keycloakService *KeycloakService) *AuthService {
	return &AuthService{
		store:           store,
		HttpClient:      httpClient,
		KeycloakService: keycloakService,
	}
}

func (service *AuthService) Login(email string, password string) (*dto.LoginDTO, error) {
	responseBody, err := service.KeycloakService.LoginKeycloakUser(email, password)
	if err != nil {
		return nil, err
	}

	loginDTO := &dto.LoginDTO{}
	json.NewDecoder(responseBody).Decode(loginDTO)
	if !service.userIsEnabled(loginDTO) {
		return nil, errors.New("User account is disabled. Check email for verification")
	}

	return loginDTO, nil
}

func (service *AuthService) SignUp(email, firstName, lastName, password, address, group string) (domain.Verification, error) {
	requestBody := dto.NewKeycloakDTO(email, firstName, lastName, password, address, group)

	loginDTO, err := service.getAdminLoginDTO()
	if err != nil {
		return domain.Verification{}, err
	}

	insertedVerificationDto := dto.VerificationDTO{}

	userId, err := service.KeycloakService.CreateKeycloakUser(requestBody, domain.BearerSchema+loginDTO.AccessToken)
	if err != nil {
		return domain.Verification{}, err
	}

	verification := &domain.Verification{
		UserId:    userId,
		FirstName: firstName,
		LastName:  lastName,
		Address:   address,
		Code:      service.generateRandomVerificationCode(),
	}

	verificationId, err := service.store.Insert(verification)
	insertedVerificationDto.Id = verificationId.Hex()
	insertedVerificationDto.UserId = userId
	if err != nil {
		return domain.Verification{}, err
	}

	return *verification, nil
}

func (service *AuthService) VerifyUser(verificationId, userId string, securityCode int) error {
	verificationPrimitiveObjectId, err := primitive.ObjectIDFromHex(verificationId)
	if err != nil {
		return err
	}
	verification, err := service.store.Get(verificationPrimitiveObjectId)
	if err != nil {
		return err
	}
	if !service.userIsCorrectForVerificationCode(verification.UserId, userId) {
		return errors.New("incorrect user for verification with id " + verificationId)
	}
	if !service.codesAreSame(verification.Code, securityCode) {

		return errors.New("verification code incorrect")
	}

	loginDTO, err := service.getAdminLoginDTO()
	if err != nil {
		return err
	}

	responseBody, err := service.KeycloakService.UpdateKeycloakUser(domain.BearerSchema+loginDTO.AccessToken, userId, dto.NewUpdateKeycloakUserDTO(userId, verification.FirstName, verification.LastName, verification.Address))

	if err != nil {
		return err
	}

	log.Println(responseBody)
	err = service.store.Delete(verificationPrimitiveObjectId)
	if err != nil {
		return err
	}

	return nil
}

func (service *AuthService) getAdminLoginDTO() (*dto.LoginDTO, error) {
	responseBody, err := service.KeycloakService.LoginKeycloakUser(domain.AdminUsername, domain.AdminPassword)
	if err != nil {
		return nil, err
	}

	loginDTO := &dto.LoginDTO{}
	json.NewDecoder(responseBody).Decode(loginDTO)
	return loginDTO, nil
}

func (service *AuthService) UpdateVerificationCode(verificationId string) (domain.Verification, error) {
	verificationPrimitiveObjectId, err := primitive.ObjectIDFromHex(verificationId)
	if err != nil {
		return domain.Verification{}, err
	}
	verification, err := service.store.Get(verificationPrimitiveObjectId)
	verification.Code = service.generateRandomVerificationCode()
	err = service.store.Update(verificationPrimitiveObjectId, verification)
	if err != nil {
		return *verification, err
	}

	return *verification, nil
}

func (service *AuthService) codesAreSame(verificationCode, inputCode int) bool {
	return verificationCode == inputCode
}

func (service *AuthService) userIsCorrectForVerificationCode(verificationUserId, authenticatedUserId string) bool {
	return verificationUserId == authenticatedUserId
}

func (service *AuthService) generateRandomVerificationCode() int {
	rand.Seed(time.Now().UnixNano())
	minRange := 1000
	maxRange := 9999
	randomNum := rand.Intn(maxRange-minRange+1) + minRange
	return randomNum
}

func (service *AuthService) userIsEnabled(loginDTO *dto.LoginDTO) bool {
	responseBody, err := service.KeycloakService.GetKeycloakUser(domain.BearerSchema + loginDTO.AccessToken)
	if err != nil {
		return false
	}

	userDTO := &dto.UserDTO{}
	json.NewDecoder(responseBody).Decode(userDTO)

	if !userDTO.EmailVerified {
		return false
	}
	return true
}
