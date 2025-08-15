package main

import (
	"collp-backend/config"
	controller "collp-backend/controllers"
	"collp-backend/middleware"
	"collp-backend/routes"
	"io/ioutil"
	"log"
	"os"

	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/unrolled/secure"
)

func SecurityMiddleware() gin.HandlerFunc {
	sec := secure.New(secure.Options{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'; img-src 'self' data:; font-src 'self' https://fonts.gstatic.com; script-src 'self' https://cdnjs.cloudflare.com",
		SSLRedirect:           false, // Set to true in production with HTTPS
		SSLProxyHeaders: map[string]string{
			"X-Forwarded-Proto": "https",
		},
	})
	return func(c *gin.Context) {
		err := sec.Process(c.Writer, c.Request)
		if err != nil {
			c.Abort()
			return
		}
		c.Next()
	}
}

var requestCounts = sync.Map{}
var lastRequest = make(map[string]time.Time)

func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if t, ok := lastRequest[ip]; ok {
			if time.Since(t) > window {
				requestCounts.Store(ip, 0)
			}
		}

		count, _ := requestCounts.LoadOrStore(ip, 0)
		if count.(int) >= limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			return
		}

		requestCounts.Store(ip, count.(int)+1)
		lastRequest[ip] = time.Now()
		c.Next()
	}
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables or defaults")
	}

	// Initialize database
	config.InitDB()

	// Initialize Google OAuth
	// Initialize auth controller
	controller.InitAuthController()

	// Load RSA private key
	keyData, err := ioutil.ReadFile("rsa.pem")
	if err != nil {
		log.Fatalf("Failed to read rsa.pem: %v", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		log.Fatalf("Failed to parse RSA private key: %v", err)
	}

	// Set public key for middleware
	middleware.SetPublicKey(&privateKey.PublicKey)

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Configure for production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Security middleware
	r.Use(SecurityMiddleware())

	// Rate limiting middleware
	r.Use(RateLimitMiddleware(100, time.Minute))

	// Setup routes
	routes.SetupRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}
