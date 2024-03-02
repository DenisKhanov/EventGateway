package services

import (
	"GoProgects/WorkTests/Exnode/FirstService/internal/producers"
	"context"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"log"
	"sync"
)

type Service struct {
	brokerAddress    string
	kafkaAdminClient sarama.ClusterAdmin
	createdTopics    map[string]struct{}
}

func NewService(brokerAddress string) *Service {
	return &Service{
		brokerAddress:    brokerAddress,
		kafkaAdminClient: *newKafkaAdminClient(brokerAddress),
		createdTopics:    make(map[string]struct{}),
	}
}

// CheckTopics проверяет существующие топики и сохраняет их в createdTopics
func (s *Service) CheckTopics() {
	topics, err := s.kafkaAdminClient.ListTopics()
	if err != nil {
		log.Fatal("Error listing topics:", err)
	}
	for k, _ := range topics {
		s.createdTopics[k] = struct{}{}
	}
}

func newKafkaAdminClient(brokerAddress string) *sarama.ClusterAdmin {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion // Указываем версию Kafka

	// Указываем адрес брокера Kafka
	client, err := sarama.NewClient([]string{brokerAddress}, config)
	if err != nil {
		log.Fatal("Error creating Kafka client:", err)
	}

	// Создаем KafkaAdminClient
	adminClient, err := sarama.NewClusterAdminFromClient(client)
	if err != nil {
		log.Fatal("Error creating Kafka admin client:", err)
	}
	return &adminClient
}

// createNewTopic создает новый топик
func (s *Service) createNewTopic(topicName string) error {

	// Указываем параметры для создания нового топика
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	// Создаем новый топик
	err := s.kafkaAdminClient.CreateTopic(topicName, topicDetail, false)
	if err != nil {
		log.Fatal("Error creating topic:", err)
		return err

	}

	logrus.Infof("Topic %s created successfully.", topicName)
	s.createdTopics[topicName] = struct{}{}
	return nil
}

// createHandlerTopic создание топика "Handler"
func (s *Service) createHandlerTopic() {
	if _, ok := s.createdTopics["Handler"]; !ok {
		if err := s.createNewTopic("Handler"); err != nil {
			logrus.Errorf("Проблема при создании топика: %v", err)
		}
	}
}

var once sync.Once

// WriteMessage проверяет наличие топика, если его нет то создает новый и записывает в него сообщение
func (s *Service) WriteMessage(ctx context.Context, topicName string, data []byte) error {
	//Создаем топик "Handler" только однажды
	once.Do(s.createHandlerTopic)
	//если топик не был создан ранее, то создаем его
	if _, ok := s.createdTopics[topicName]; !ok {
		if err := s.createNewTopic(topicName); err != nil {
			logrus.Errorf("Проблема при создании топика: %v", err)
		}
	}
	//запускаем создание продюсера и запись сообщения в топик
	if err := producers.ProduceMessage(s.brokerAddress, topicName, data); err != nil {
		logrus.Errorf("Проблема при записи в топик: %v", err)
		return err
	}
	return nil
}

// WriteMessageSarama проверяет наличие топика, если его нет то создает новый и записывает в него сообщение
func (s *Service) WriteMessageSarama(topicName string, data []byte) error {
	//запускаем создание продюсера и запись сообщения в топик
	if err := producers.ProduceMessage(s.brokerAddress, topicName, data); err != nil {
		logrus.Errorf("Проблема при записи в топик: %v", err)
		return err
	}
	return nil
}
