package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mmmajder/zms-devops-auth-service/application"
	"github.com/mmmajder/zms-devops-auth-service/domain"
	"github.com/mmmajder/zms-devops-auth-service/infrastructure/request"
	"log"
	"net/http"
)

type UserHandler struct {
	service *application.UserService
}

func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (handler *UserHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", domain.UserContextPath, handler.GetUser)
	err = mux.HandlePath("GET", domain.UserContextPath+"/{id}", handler.GetUserById)
	err = mux.HandlePath("PUT", domain.UserContextPath+"/{id}", handler.UpdateUser)
	err = mux.HandlePath("PUT", domain.UserContextPath+"/{id}/reset-password", handler.ResetPassword)
	err = mux.HandlePath("DELETE", domain.UserContextPath+"/{id}", handler.DeleteUser)
	if err != nil {
		panic(err)
	}
}

func (handler *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, params map[string]string) {
	res, err := handler.service.GetUser(r.Header.Get("Authorization"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (handler *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id := params["id"]
	if id == "" {
		writeErrorMessageIfParamIsEmpty(w, "Must specify user id when updating user data")
		return
	}
	updatingUserRequest := new(request.UpdatingUserRequest)
	err := json.NewDecoder(r.Body).Decode(updatingUserRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid updating user data"))
		return
	}
	validate := validator.New()
	err = validate.Struct(updatingUserRequest)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Validation error: %s", err.(validator.ValidationErrors)),
			http.StatusBadRequest,
		)
		return
	}

	err = handler.service.UpdateUser(
		r.Header.Get("Authorization"),
		id,
		updatingUserRequest.FirstName,
		updatingUserRequest.LastName,
		updatingUserRequest.Address,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	log.Println("Updating user account data success")
	w.WriteHeader(http.StatusOK)
}

func (handler *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id := params["id"]
	if id == "" {
		writeErrorMessageIfParamIsEmpty(w, "Must specify user id when deleting user")
		return
	}
	err := handler.service.DeleteUser(r.Header.Get("Authorization"), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	log.Printf("Deleting user account success")
	w.WriteHeader(http.StatusOK)
}

func (handler *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id := params["id"]
	if id == "" {
		writeErrorMessageIfParamIsEmpty(w, "Must specify user id when resetting password")
		return
	}

	updatePasswordRequest := new(request.UpdatePasswordRequest)
	err := json.NewDecoder(r.Body).Decode(updatePasswordRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid resetting user data"))
		return
	}
	validate := validator.New()
	err = validate.Struct(updatePasswordRequest)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Validation error: %s", err.(validator.ValidationErrors)),
			http.StatusBadRequest,
		)
		return
	}

	err = updatePasswordRequest.ArePasswordsMatch()
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf(err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	err = handler.service.ResetPassword(r.Header.Get("Authorization"), id, updatePasswordRequest.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	log.Printf("Resetting password success")
	w.WriteHeader(http.StatusOK)
}

func (handler *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request, params map[string]string) {
	id := params["id"]
	if id == "" {
		writeErrorMessageIfParamIsEmpty(w, "Must specify user id when getting user by id")
		return
	}

	res, err := handler.service.GetUserById(r.Header.Get("Authorization"), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
