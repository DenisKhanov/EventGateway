package handlers

import (
	"GoProgects/WorkTests/Exnode/FirstService/internal/auth"
	"GoProgects/WorkTests/Exnode/FirstService/internal/models"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type Service interface {
	WriteMessage(ctx context.Context, topicName string, data []byte) error
	WriteMessageSarama(topicName string, data []byte) error
	CheckTopics()
}

type Handler struct {
	Service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h Handler) GetTopicName(c *gin.Context) {
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
	if err = h.Service.WriteMessageSarama(topicName, body); err != nil {
		c.JSON(http.StatusGatewayTimeout, gin.H{"Error": err})
		return
	}
	//if err = h.Service.WriteMessage(ctx, topicName, body); err != nil {
	//	c.JSON(http.StatusGatewayTimeout, gin.H{"Error": err})
	//	return
	//}
}

// MiddlewareLogging provides a logging middleware for Gin.
// It logs details about each request including the URL, method, response status, duration, and size.
// This middleware is useful for monitoring and debugging purposes.
func (h Handler) MiddlewareLogging() gin.HandlerFunc {
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
func (h Handler) MiddlewareAuthPrivate() gin.HandlerFunc {
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
