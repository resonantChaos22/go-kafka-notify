package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/resonantchaos22/go-kafka-notify/pkg/models"
)

const (
	ProducerPort       = ":8080"
	KafkaServerAddress = "localhost:9092"
	KafkaTopic         = "notifications"
)

var ErrUserNotFound = errors.New("user not found")

type Producer struct {
	syncProducer sarama.SyncProducer
	users        []models.User
}

// KAFKA FUNCTIONS
func NewProducer(users []models.User) (*Producer, error) {
	producer := new(Producer)
	err := producer.setup(users)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}
	return producer, nil
}

func (p *Producer) setup(users []models.User) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{KafkaServerAddress}, config)

	if err != nil {
		return fmt.Errorf("failed to setup producer: %w", err)
	}
	p.syncProducer = producer
	p.users = users

	return nil
}

func (p *Producer) sendKafkaMessage(ctx *gin.Context, fromID, toID int) error {
	message := ctx.PostForm("message")

	fromUser, err := p.findByUserId(fromID)
	if err != nil {
		return err
	}

	toUser, err := p.findByUserId(toID)
	if err != nil {
		return err
	}

	notification := models.Notification{
		From:    fromUser,
		To:      toUser,
		Message: message,
	}

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: KafkaTopic,
		Key:   sarama.StringEncoder(strconv.Itoa(toUser.ID)),
		Value: sarama.StringEncoder(notificationJSON),
	}

	_, _, err = p.syncProducer.SendMessage(msg)
	return err
}

func (p *Producer) sendMessageHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fromID, err := p.getIDFromRequest("fromID", ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		toID, err := p.getIDFromRequest("toID", ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		err = p.sendKafkaMessage(ctx, fromID, toID)
		if errors.Is(err, ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "User Not Found"})
			return
		}

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Notification sent successfully!",
		})
	}
}

//	HELPER FUNCTIONS

// findByUserId is a helper finds a user with a given id in a list of users.
func (p *Producer) findByUserId(id int) (models.User, error) {
	for _, user := range p.users {
		if user.ID == id {
			return user, nil
		}
	}

	return models.User{}, ErrUserNotFound
}

// getIDFromRequest parses id from the request based on the 'formValue' field
func (p *Producer) getIDFromRequest(formValue string, ctx *gin.Context) (int, error) {
	id, err := strconv.Atoi(ctx.PostForm(formValue))
	if err != nil {
		return 0, fmt.Errorf("failed to parse ID from form value %s: %w", formValue, err)
	}

	return id, nil
}

func main() {
	users := []models.User{
		{ID: 1, Name: "Shreyash"},
		{ID: 2, Name: "Rajesh"},
		{ID: 3, Name: "Riya"},
		{ID: 4, Name: "Sally"},
	}

	producer, err := NewProducer(users)
	if err != nil {
		log.Fatalf("failed to initialize producer: %v", err)
	}
	defer producer.syncProducer.Close()

	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.POST("/send", producer.sendMessageHandler())

	log.Printf("Kafka Producer started at http://localhost%s\n", ProducerPort)

	if err := router.Run(ProducerPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}

}
