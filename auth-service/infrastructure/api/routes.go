package api

//import (
//	"encoding/json"
//	"github.com/mmmajder/devops-booking-service/domain"
//	"net/http"
//)

type GreetRes struct {
	Hello string `json:"hello"`
}

//func (handler *Server) handleGreet(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodGet {
//		w.WriteHeader(400)
//		w.Write([]byte("Invalid method"))
//		return
//	}
//
//	w.Header().Add("Content-Type", "application/json")
//	w.WriteHeader(200)
//	res := &GreetRes{
//		Hello: "Hello",
//	}
//	json.NewEncoder(w).Encode(res)
//}

//func (handler *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodPost {
//		w.WriteHeader(400)
//		w.Write([]byte("Invalid method"))
//		return
//	}
//
//	loginDTO := new(domain.LoginDTO)
//	err := json.NewDecoder(r.Body).Decode(loginDTO)
//	if err != nil {
//		w.WriteHeader(400)
//		w.Write([]byte("Invalid login data"))
//		return
//	}
//
//	loginPayload := &domain.LoginPayload{
//		ClientId:  "Istio",
//		Username:  loginDTO.Username,
//		Password:  loginDTO.Password,
//		GrantType: "password",
//	}
//
//	res, err := handler.client.Login(loginPayload)
//	if err != nil {
//		w.WriteHeader(400)
//		w.Write([]byte(err.Error()))
//		return
//	}
//	w.Header().Add("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	json.NewEncoder(w).Encode(res)
//}
