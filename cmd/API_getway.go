package main

import (
	"GoProgects/WorkTests/Exnode/FirstService/internal/auth"
	"GoProgects/WorkTests/Exnode/FirstService/internal/configs"
	"GoProgects/WorkTests/Exnode/FirstService/internal/handlers"
	"GoProgects/WorkTests/Exnode/FirstService/internal/loggers"
	"GoProgects/WorkTests/Exnode/FirstService/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {

	//Инициализируем переменные окружения или флаги CLI
	cfg := configs.NewConfig()

	loggers.RunLoggerConfig(cfg.EnvLogsLevel)
	logrus.Infof("BOT started with configuration logs level: %v", cfg.EnvLogsLevel)

	myService := services.NewService(cfg.EnvBrokerAddress)
	//TODO подумать более оптимальной реализацией CheckTopics
	myService.CheckTopics()
	myHandler := handlers.NewHandler(myService)

	token, err := auth.BuildJWTString("Test_topic")
	if err != nil {
		logrus.Error(err)
	}
	logrus.Infof(token)
	// Установка переменной окружения для включения режима разработки
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	//Private middleware routers group
	privateRoutes := router.Group("/")
	privateRoutes.Use(myHandler.MiddlewareAuthPrivate())
	privateRoutes.Use(myHandler.MiddlewareLogging())
	//TODO подумать нужно ли использовать сжатие
	//privateRoutes.Use(myHandler.MiddlewareCompress())
	privateRoutes.POST("/test", myHandler.GetTopicName)
	server := &http.Server{Addr: cfg.EnvServAdr, Handler: router}
	logrus.Info("Starting server on: ", cfg.EnvServAdr)

	if err = server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logrus.Error(err)
	}
}
