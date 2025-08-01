package router
import (
	"collp-backend/controller"
	"collp-backend/middleware"
	"github.com/gorilla/mux"
)

func ThemesRouter() *mux.Router {
	publicRouters := mux.NewRouter()
	// Routes public Authentication with google
	publicRouters.HandleFunc("/auth/google/login", controller.GoogleLogin).Methods("GET")
	publicRouters.HandleFunc("/auth/google/callback", controller.GoogleCallback).Methods("GET")
	publicRouters.HandleFunc("/api/collp/login", controller.CollPLogin).Methods("POST")
	publicRouters.HandleFunc("/api/collp/register", controller.CollPRegister).Methods("POST")
	// Create private Subrouter router for auth
	privateRoutes := publicRouters.PathPrefix("/").Subrouter()
	privateRoutes.Use(middleware.AuthMiddleware)
	privateRoutes.HandleFunc("/api/collp/main-menu", controller.MainMenu).Methods("GET")

	return publicRouters
}
