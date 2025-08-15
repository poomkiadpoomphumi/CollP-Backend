package controllers

import (
	"collp-backend/services"
	"encoding/json"
	"net/http"
)

func MainMenu(w http.ResponseWriter, r *http.Request) {
	menu := services.GetAllMainMenu()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menu)
}
