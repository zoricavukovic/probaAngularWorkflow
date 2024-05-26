package api

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

type Handler interface {
	Init(mux *runtime.ServeMux)
}

func writeErrorMessageIfParamIsEmpty(w http.ResponseWriter, errorMessage string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(errorMessage))
}
