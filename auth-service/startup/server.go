package api

import (
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mmmajder/zms-devops-auth-service/application"
	"github.com/mmmajder/zms-devops-auth-service/domain"
	"github.com/mmmajder/zms-devops-auth-service/infrastructure/api"
	"github.com/mmmajder/zms-devops-auth-service/infrastructure/persistence"
	"github.com/mmmajder/zms-devops-auth-service/startup/config"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
)

type Server struct {
	config *config.Config
	mux    *runtime.ServeMux
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
		mux:    runtime.NewServeMux(),
	}
}

func (server *Server) Start() {
	server.setupHandlers()
	fmt.Println("HEJ6")
	http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), server.mux)
	log.SetOutput(os.Stdin)
	log.SetOutput(os.Stderr)
	log.SetOutput(os.Stdout)
}

func (server *Server) setupHandlers() {
	keycloakService := server.initKeycloakService()
	emailService := server.initEmailService()

	mongoClient := server.initMongoClient()
	verificationStore := server.initVerificationStore(mongoClient)
	authService := server.initAuthService(verificationStore, keycloakService)
	authHandler := server.initAuthHandler(authService, emailService)
	authHandler.Init(server.mux)

	userService := server.initUserService(authService, keycloakService)
	userHandler := server.initUserHandler(userService)
	userHandler.Init(server.mux)
}

func (server *Server) initKeycloakService() *application.KeycloakService {
	return application.NewKeycloakService(&http.Client{}, server.config.IdentityProviderHost)
}

func (server *Server) initEmailService() *application.EmailService {
	return application.NewEmailService()
}

func (server *Server) initAuthService(store domain.VerificationStore, keycloakService *application.KeycloakService) *application.AuthService {
	return application.NewAuthService(store, &http.Client{}, keycloakService)
}

func (server *Server) initAuthHandler(authService *application.AuthService, emailService *application.EmailService) *api.AuthHandler {
	return api.NewAuthHandler(authService, emailService)
}

func (server *Server) initUserService(authService *application.AuthService, keycloakService *application.KeycloakService) *application.UserService {
	return application.NewUserService(&http.Client{}, authService, keycloakService, server.config.IdentityProviderHost)
}

func (server *Server) initUserHandler(service *application.UserService) *api.UserHandler {
	return api.NewUserHandler(service)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.DBUsername, server.config.DBPassword, server.config.DBHost, server.config.DBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initVerificationStore(client *mongo.Client) domain.VerificationStore {
	return persistence.NewVerificationMongoDBStore(client)
}
