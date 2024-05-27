package main

import (
	"email-service/config"
	"email-service/model"
	"email-service/service"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	channel, close := config.ConnectRabbitmqs()
	defer close()
	defer channel.Close()
	emailConfig := config.NewEmailConfig()
	emailService := service.NewEmailService(emailConfig)
	fmt.Println("Check")
	go func() {
		fmt.Println("Check2")
		emailComsumer, err := channel.Consume("email", "consumer-email", true, false, false, false, nil)
		fmt.Println("Check3")
		if err != nil {
			fmt.Println(err)
		}
		for message := range emailComsumer {
			var response model.EmailRequest
			if err := json.Unmarshal(message.Body, &response); err != nil {
				message.Nack(false, false)
				log.Printf("failed to unmarshal order: %v", err)
				continue
			}
			emailService.SendEmail(response.Subject, response.UserEmail, response.Message)
		}
	}()
	r := gin.Default()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
