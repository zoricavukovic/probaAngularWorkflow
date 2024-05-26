package application

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mmmajder/zms-devops-auth-service/domain"
	"github.com/mmmajder/zms-devops-auth-service/infrastructure/dto"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type KeycloakService struct {
	HttpClient           *http.Client
	IdentityProviderHost string
}

func NewKeycloakService(httpClient *http.Client, identityProviderHost string) *KeycloakService {

	return &KeycloakService{
		HttpClient:           httpClient,
		IdentityProviderHost: identityProviderHost,
	}
}

func (service *KeycloakService) CreateKeycloakUser(signupDTO *dto.KeycloakDTO, authorizationHeader string) (string, error) {
	jsonBody, err := json.Marshal(signupDTO)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://%s/admin/realms/Istio/users", service.IdentityProviderHost),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return "", err
	}
	setContentType(req, domain.JsonContentType)
	req.Header.Set(domain.Authorization, authorizationHeader)

	resp, err := service.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusConflict {
			return "", errors.New("user exists with same email")
		}
		return "", errors.New("registration failed")
	}

	location := resp.Header.Get("Location")
	if location != "" {
		parts := strings.Split(location, "/")
		if len(parts) > 0 {
			lastPart := parts[len(parts)-1]
			return lastPart, nil
		} else {
			return "", errors.New("user not created")
		}
	} else {
		return "", errors.New("user not created")
	}
}

func setContentType(req *http.Request, contentType string) {
	req.Header.Set(domain.ContentType, contentType)
}

func setRequestHeader(req *http.Request, contentType string) {
	req.Header.Set(domain.ContentType, contentType)
	req.Header.Add("Host", "keycloak.backend.svc.cluster.local")
}

func (service *KeycloakService) LoginKeycloakUser(email, password string) (io.ReadCloser, error) {
	formData := url.Values{
		"client_id":  {"Istio"},
		"username":   {email},
		"password":   {password},
		"grant_type": {"password"},
		"scope":      {"openid"},
	}
	encodedData := formData.Encode()

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://keycloak.backend.svc.cluster.local/realms/Istio/protocol/openid-connect/token"),
		strings.NewReader(encodedData),
	)
	if err != nil {
		return nil, err
	}

	setRequestHeader(req, "application/x-www-form-urlencoded")

	resp, err := service.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("cao login 3")
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("login failed")
	}

	return resp.Body, nil
}

func (service *KeycloakService) GetKeycloakUser(authorizationHeader string) (io.ReadCloser, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://%s/realms/Istio/protocol/openid-connect/userinfo", service.IdentityProviderHost),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set(domain.Authorization, authorizationHeader)

	resp, err := service.HttpClient.Do(req)
	log.Printf("---------")
	log.Println(resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("getting user failed")
	}

	return resp.Body, nil
}

func (service *KeycloakService) GetKeycloakUserById(authorizationHeader, id string) (io.ReadCloser, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://%s/admin/realms/Istio/users/%s?q=address", service.IdentityProviderHost, id),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set(domain.Authorization, authorizationHeader)

	resp, err := service.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("getting user failed")
	}

	return resp.Body, nil
}

func (service *KeycloakService) UpdateKeycloakUser(authorizationHeader string, id string, updateUser *dto.UpdateKeycloakUserDTO) (io.ReadCloser, error) {
	jsonBody, err := json.Marshal(updateUser)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("http://%s/admin/realms/Istio/users/%s", service.IdentityProviderHost, id),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}
	setContentType(req, domain.JsonContentType)
	req.Header.Set(domain.Authorization, authorizationHeader)

	resp, err := service.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusNoContent {
		return nil, errors.New("updating user failed")
	}

	return resp.Body, nil
}

func (service *KeycloakService) DeleteKeycloakUser(authorizationHeader, id string) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("http://%s/admin/realms/Istio/users/%s", service.IdentityProviderHost, id),
		nil,
	)
	if err != nil {
		return err
	}
	req.Header.Set(domain.Authorization, authorizationHeader)

	resp, err := service.HttpClient.Do(req)
	log.Println("Deleting user")
	log.Printf(resp.Status)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.New("deleting user failed")
	}
	return nil
}

func (service *KeycloakService) ResetPasswordOfKeycloakUser(authorizationHeader string, id string, requestBody *dto.CredentialsDTO) error {
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("http://%s/admin/realms/Istio/users/%s/reset-password", service.IdentityProviderHost, id),
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return err
	}
	setContentType(req, domain.JsonContentType)
	req.Header.Set(domain.Authorization, authorizationHeader)

	resp, err := service.HttpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.New("resetting user password failed")
	}
	return nil
}
