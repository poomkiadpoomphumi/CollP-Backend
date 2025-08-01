package controller
import (
	"collp-backend/usecase"
	"encoding/json"
	"net/http"
)

func CollPLogin(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
    password := r.Header.Get("password")
	response := usecase.CollPLoginUsecase(username, password)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CollPRegister(w http.ResponseWriter, r *http.Request) {
    username := r.Header.Get("username")
    password := r.Header.Get("password")
    phone := r.Header.Get("phone")
    address := r.Header.Get("address")
    response := usecase.CollPRegisterUsecase(username, password, phone, address)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
