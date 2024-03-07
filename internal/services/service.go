package services

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"sync"
)

type Consumer interface {
	ListenKafkaMessages(topicName string, message chan<- *sarama.ConsumerMessage) error
}
type Producer interface {
	SandMsgToKafka(ctx context.Context, topicName, requestID string, data []byte) error
}

type ServiceKafka struct {
	Producer         Producer
	Customer         Consumer
	ResponseChannels map[string]chan *sarama.ConsumerMessage
	ConsumerState    map[string]struct{}
	message          chan *sarama.ConsumerMessage
	Mu               sync.Mutex
	Once             sync.Once
}

func NewServiceKafka(producer Producer, consumer Consumer) *ServiceKafka {
	return &ServiceKafka{
		Producer:         producer,
		Customer:         consumer,
		ResponseChannels: make(map[string]chan *sarama.ConsumerMessage),
		ConsumerState:    make(map[string]struct{}),
		message:          make(chan *sarama.ConsumerMessage),
	}
}

// SendMsgToKafka sends a message to the specified Kafka topic using the underlying producer.
// It initiates the producer and writes the message to the topic.
func (s *ServiceKafka) SendMsgToKafka(ctx context.Context, topicName, requestID string, data []byte) error {
	//запускаем создание продюсера и запись сообщения в топик
	if err := s.Producer.SandMsgToKafka(ctx, topicName, requestID, data); err != nil {
		logrus.Errorf("Проблема при записи в топик: %v", err)
		return err
	}
	return nil
}

// RunTopicConsumer starts a Kafka topic consumer for the given topic.
// It ensures that the consumer is initialized only once and starts listening for messages on the specified topic.
func (s *ServiceKafka) RunTopicConsumer(ctx context.Context, topicName string) {
	go s.Once.Do(s.runListenerConsumers)
	_, exists := s.ConsumerState[topicName]
	if !exists {
		s.ConsumerState[topicName] = struct{}{}
		// Запуск нового консьюмера
		errCh := make(chan error, 1)
		go func() {
			logrus.Infof("Create Consumer %s", topicName)
			errCh <- s.Customer.ListenKafkaMessages(topicName, s.message)
		}()

		select {
		case <-ctx.Done():
			return
		case err := <-errCh:
			if err != nil {
				delete(s.ConsumerState, topicName)
				logrus.Infof("Consumer %s deleted", topicName)
				return
			}
		}
	}

}

// runListenerConsumers continuously listens for messages on the message channel.
// It processes incoming messages, sends them to the appropriate response channel, and handles channel closure.
func (s *ServiceKafka) runListenerConsumers() {
	logrus.Info("ListenerConsumers is started")
	for msg := range s.message {
		responseID := string(msg.Key)
		s.Mu.Lock()
		ch, exists := s.ResponseChannels[responseID]
		if exists {
			ch <- msg
			delete(s.ResponseChannels, responseID)
			logrus.Infof("RunListener download %s", responseID)
		}
		s.Mu.Unlock()
	}
}
