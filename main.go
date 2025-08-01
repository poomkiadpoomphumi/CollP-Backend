package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	/*
		"strconv"

		"collp-backend/connection"
	*/
	"collp-backend/controller"
	"collp-backend/middleware"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func main() {
	// โหลดไฟล์ .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables or defaults")
	}
	controller.InitGoogleOauth()
	// โหลด private key จากไฟล์ PEM
	keyData, err := ioutil.ReadFile("rsa.pem")
	if err != nil {
		log.Fatalf("Failed to read rsa.pem: %v", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		log.Fatalf("Failed to parse RSA private key: %v", err)
	}
	// ตั้ง public key ให้ middleware ใช้ verify token
	middleware.SetPublicKey(&privateKey.PublicKey)
	// อ่านค่า config จาก env
	/* 	user := os.Getenv("ORACLE_USER")
	   	pass := os.Getenv("ORACLE_PASS")
	   	host := os.Getenv("ORACLE_HOST")
	   	portStr := os.Getenv("ORACLE_PORT")
	   	service := os.Getenv("ORACLE_SERVICE")

	   	if user == "" || pass == "" || host == "" || portStr == "" || service == "" {
	   		log.Fatal("Missing one or more required environment variables for Oracle connection")
	   	}

	   	port, err := strconv.Atoi(portStr)
	   	if err != nil {
	   		log.Fatalf("Invalid ORACLE_PORT: %v", err)
	   	}

	   	// สร้าง OracleDB connection ให้ถูกต้อง
	   	oracle := connection.NewOracleDB(user, pass, host, port, service)
	   	defer oracle.Close() */

	// Run server
	// Read port from environment or fallback to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := NewServer()
	log.Printf("Server running at port :%s", port)
	server.Run(fmt.Sprintf(":%s", port))
}
