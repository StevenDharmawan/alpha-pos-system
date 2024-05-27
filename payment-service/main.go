package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"payment-service/config"
	"payment-service/controller"
	"payment-service/model"
	"payment-service/service"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	channel, close := config.ConnectRabbitmq()
	defer close()
	defer channel.Close()

	midtransService := service.NewMidtransService()
	go func() {
		paymentComsumer, err := channel.Consume("payment-request", "consumer-payment", true, false, false, false, nil)
		if err != nil {
			log.Fatal(err)
		}
		for message := range paymentComsumer {
			var payment model.PaymentRequest
			if err := json.Unmarshal(message.Body, &payment); err != nil {
				message.Nack(false, false)
				log.Printf("failed to unmarshal order: %v", err)
				continue
			}
			url, err := midtransService.GenerateSnapURL(payment)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(url)
			marshalledUrl, _ := json.Marshal(model.EmailRequest{
				Subject:   "Payment Url",
				UserEmail: "stevenixe123@gmail.com",
				Message:   url,
			})
			err = channel.PublishWithContext(context.Background(), "email", "email", false, false, amqp.Publishing{ContentType: "application/json", Body: marshalledUrl})
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
	midtransController := controller.NewMidtransController(midtransService, channel)

	r := gin.Default()
	r.POST("midtrans/call-back", midtransController.PaymentHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
