package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mmmajder/zms-devops-auth-service/application"
	"github.com/mmmajder/zms-devops-auth-service/domain"
	"github.com/mmmajder/zms-devops-auth-service/infrastructure/dto"
	"github.com/mmmajder/zms-devops-auth-service/infrastructure/request"
	"log"
	"net/http"
)

type AuthHandler struct {
	authService  *application.AuthService
	emailService *application.EmailService
}

func NewAuthHandler(authService *application.AuthService, emailService *application.EmailService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		emailService: emailService,
	}
}

type EmailData struct {
	verificationCode string
}

func (handler *AuthHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("POST", domain.AuthContextPath+"/login", handler.Login)
	err = mux.HandlePath("POST", domain.AuthContextPath+"/signup", handler.SignUp)
	err = mux.HandlePath("PUT", domain.AuthContextPath+"/verify", handler.VerifyUser)
	err = mux.HandlePath("POST", domain.AuthContextPath+"/send-code-again", handler.SendVerificationCodeAgain)
	err = mux.HandlePath("GET", domain.AuthContextPath+"/health", handler.GetHealthCheck)
	if err != nil {
		panic(err)
	}
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	loginRequest := new(request.LoginRequest)
	err := json.NewDecoder(r.Body).Decode(loginRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid login data"))
		return
	}
	validate := validator.New()
	err = validate.Struct(loginRequest)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Validation error: %s", err.(validator.ValidationErrors)),
			http.StatusBadRequest,
		)
		return
	}

	res, err := handler.authService.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (handler *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	registrationRequest := new(request.RegistrationRequest)
	err := json.NewDecoder(r.Body).Decode(registrationRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid registration data"))
		return
	}
	validate := validator.New()
	err = validate.Struct(registrationRequest)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Validation error: %s", err.(validator.ValidationErrors)),
			http.StatusBadRequest,
		)
		return
	}

	err = registrationRequest.ArePasswordsMatch()
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf(err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	verification, err := handler.authService.SignUp(
		registrationRequest.Email,
		registrationRequest.FirstName,
		registrationRequest.LastName,
		registrationRequest.Password,
		registrationRequest.Address,
		registrationRequest.Group,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	handler.emailService.SendEmail(domain.SubjectVerifyUser, handler.emailService.GetSendVerificationCodeEmailBody(registrationRequest.Email, verification))

	log.Printf("SignUp success")

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.VerificationDTO{Id: verification.Id.Hex(), UserId: verification.UserId})
}

func (handler *AuthHandler) VerifyUser(w http.ResponseWriter, r *http.Request, params map[string]string) {
	verificationRequest := new(request.VerificationRequest)
	err := json.NewDecoder(r.Body).Decode(verificationRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid verification data"))
		return
	}
	validate := validator.New()
	err = validate.Struct(verificationRequest)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Validation error: %s", err.(validator.ValidationErrors)),
			http.StatusBadRequest,
		)
		return
	}

	err = handler.authService.VerifyUser(
		verificationRequest.VerificationId,
		verificationRequest.UserId,
		verificationRequest.SecurityCode,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	log.Printf("Verify success")
	w.WriteHeader(http.StatusOK)
}

func (handler *AuthHandler) SendVerificationCodeAgain(w http.ResponseWriter, r *http.Request, params map[string]string) {
	sendCodeAgainRequest := new(request.SendVerificationCodeRequest)
	err := json.NewDecoder(r.Body).Decode(sendCodeAgainRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid data for send new verification code method"))
		return
	}
	validate := validator.New()
	err = validate.Struct(sendCodeAgainRequest)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Validation error: %s", err.(validator.ValidationErrors)),
			http.StatusBadRequest,
		)
		return
	}

	verification, err := handler.authService.UpdateVerificationCode(
		sendCodeAgainRequest.VerificationId,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	verificationDTO := dto.VerificationDTO{Id: verification.Id.Hex(), UserId: verification.UserId}

	handler.emailService.SendEmail(domain.SubjectVerifyUser, handler.emailService.GetSendVerificationCodeEmailBody(sendCodeAgainRequest.UserEmail, verification))

	log.Printf("Send code again success")
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(verificationDTO)
}

func (handler *AuthHandler) GetHealthCheck(w http.ResponseWriter, r *http.Request, params map[string]string) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("AUTH SERVICE IS HEALTH")
}
