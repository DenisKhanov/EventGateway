package main

import (
	"GoProgects/WorkTests/Exnode/FirstService/internal/configs"
	"GoProgects/WorkTests/Exnode/FirstService/internal/consumers"
	"GoProgects/WorkTests/Exnode/FirstService/internal/handlers"
	"GoProgects/WorkTests/Exnode/FirstService/internal/loggers"
	"GoProgects/WorkTests/Exnode/FirstService/internal/producers"
	"GoProgects/WorkTests/Exnode/FirstService/internal/services"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	//Инициализируем переменные окружения или флаги CLI
	cfg := configs.NewConfig()

	loggers.RunLoggerConfig(cfg.EnvLogsLevel)
	logrus.Infof("First service started with configuration logs level: %v", cfg.EnvLogsLevel)

	myProducer, err := producers.NewProducerKafka(cfg.EnvBrokerAddress)
	if err != nil {
		logrus.Fatalf("Failed to create producer: %v", err)
	}
	defer myProducer.Close()
	myConsumer, err := consumers.NewConsumerKafka(cfg.EnvBrokerAddress)
	if err != nil {
		logrus.Fatalf("Failed to create consumer: %v", err)
	}
	defer myConsumer.Close()

	myService := services.NewServiceKafka(myProducer, myConsumer)

	myHandler := handlers.NewHandlerKafka(myService, myService)

	// Установка переменной окружения для включения режима разработки
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	router.GET("/token", myHandler.GetTokens)
	//Private middleware routers group
	privateRoutes := router.Group("/")
	privateRoutes.Use(myHandler.MiddlewareAuthPrivate())
	privateRoutes.Use(myHandler.MiddlewareLogging())

	//TODO подумать нужно ли использовать сжатие
	//privateRoutes.Use(myHandler.MiddlewareCompress())

	privateRoutes.POST("/test", myHandler.ChangeJSONData)
	server := &http.Server{
		Addr:         cfg.EnvServAdr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logrus.Info("Starting server on: ", cfg.EnvServAdr)

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logrus.Error(err)
		}
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-signalChan:
		logrus.Infof("Shutting down server with signal : %v...", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = server.Shutdown(ctx); err != nil {
			logrus.Errorf("HTTP server Shutdown error: %v\n", err)
		}
	}
	logrus.Info("Server exited")
}
