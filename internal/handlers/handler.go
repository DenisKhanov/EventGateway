package handlers

import (
	"GoProgects/WorkTests/Exnode/FirstService/internal/auth"
	"GoProgects/WorkTests/Exnode/FirstService/internal/models"
	"GoProgects/WorkTests/Exnode/FirstService/internal/services"
	"context"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type Service interface {
	SendMsgToKafka(ctx context.Context, topicName, requestID string, data []byte) error
	RunTopicConsumer(ctx context.Context, topicName string)
}

type HandlerKafka struct {
	ServiceKafka *services.ServiceKafka
	Service      Service
}

func NewHandlerKafka(ser *services.ServiceKafka, service Service) *HandlerKafka {
	return &HandlerKafka{
		ServiceKafka: ser,
		Service:      service}
}

// GetTokens generates and returns JWT tokens corresponding to predefined topic names.
// It is used for authentication and obtaining topic-related tokens.
func (h HandlerKafka) GetTokens(c *gin.Context) {
	topicNames := []string{
		"One",
		"Two",
		"Three",
		"Four",
		"Five",
		"Six",
		"Seven",
		"Eight",
		"Nine",
		"Ten",
	}
	var tokens []string
	for _, v := range topicNames {
		token, err := auth.BuildJWTString(v)
		if err != nil {
			logrus.Error(err)
		}
		tokens = append(tokens, token)
	}
	c.JSON(http.StatusOK, tokens)
	return
}

// ChangeJSONData handles incoming requests to modify JSON data, sending it to a Kafka topic
// and asynchronously listening for the processed message. It returns the modified message or a timeout error.
func (h HandlerKafka) ChangeJSONData(c *gin.Context) {
	ctx := c.Request.Context()
	topicName, ok := ctx.Value(models.TopicName).(string)
	if !ok {
		logrus.Fatalf("context value is not topic name: %s", topicName)
	}
	// Читаем JSON тело запроса в виде среза байтов
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading request body"})
		return
	}
	// создаем уникальный ID запроса, для соблюдения соответствия запроса и возвращаемых из Kafka сообщений
	requestID := uuid.New().String()
	responseCh := make(chan *sarama.ConsumerMessage)
	h.ServiceKafka.Mu.Lock()
	h.ServiceKafka.ResponseChannels[requestID] = responseCh
	h.ServiceKafka.Mu.Unlock()

	go h.ServiceKafka.RunTopicConsumer(ctx, topicName)

	// Запускаем отправку сообщения в Kafka
	if err = h.ServiceKafka.SendMsgToKafka(ctx, topicName, requestID, body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}

	select {
	case responseMsg := <-responseCh:
		c.JSON(200, gin.H{"message": string(responseMsg.Value)})
	case <-time.After(10 * time.Second):
		h.ServiceKafka.Mu.Lock()
		delete(h.ServiceKafka.ResponseChannels, requestID)
		h.ServiceKafka.Mu.Unlock()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "timeout waiting for response"})
	}

}

// MiddlewareLogging provides a logging middleware for Gin.
// It logs details about each request including the URL, method, response status, duration, and size.
// This middleware is useful for monitoring and debugging purposes.
func (h HandlerKafka) MiddlewareLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Запуск таймера
		start := time.Now()

		// Обработка запроса
		c.Next()

		// Измерение времени обработки
		duration := time.Since(start)

		// Получение статуса ответа и размера
		status := c.Writer.Status()
		size := c.Writer.Size()

		// Логирование информации о запросе
		logrus.WithFields(logrus.Fields{
			"url":      c.Request.URL.RequestURI(),
			"method":   c.Request.Method,
			"status":   status,
			"duration": duration,
			"size":     size,
		}).Info("Обработан запрос")
	}
}

// MiddlewareAuthPrivate provides authentication middleware for private routes.
// It checks the user token and only allows access if the token is valid.
// This middleware ensures that only authenticated users can access certain routes.
func (h HandlerKafka) MiddlewareAuthPrivate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Не предоставлен токен авторизации"})
		}
		token = strings.TrimPrefix(token, "Bearer ")
		topicName, err := auth.GetUserTopicName(token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		ctx := context.WithValue(c.Request.Context(), models.TopicName, topicName)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
