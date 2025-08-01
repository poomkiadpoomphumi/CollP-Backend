package controller
import (
	"collp-backend/usecase"
	"encoding/json"
	"net/http"
)

func MainMenu(w http.ResponseWriter, r *http.Request) {
	menu := usecase.GetAllMainMenu()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menu)
}