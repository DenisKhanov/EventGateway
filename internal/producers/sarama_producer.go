package producers

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type ProducerKafka struct {
	brokerAddress  string
	saramaProducer sarama.SyncProducer
}

func NewProducerKafka(brokerAddress string) (*ProducerKafka, error) {
	saramaProducer, err := createKafkaProducer(brokerAddress)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &ProducerKafka{
		brokerAddress:  brokerAddress,
		saramaProducer: saramaProducer,
	}, nil

}

// createKafkaProducer  creates a synchronous Kafka producer.
// It takes the broker address as a parameter and returns the producer and any error encountered during creation.
func createKafkaProducer(brokerAddress string) (sarama.SyncProducer, error) {
	producer, err := sarama.NewSyncProducer([]string{brokerAddress}, nil)
	if err != nil {
		logrus.Errorf("Failed to create producer: %v", err)
		return nil, err
	}
	return producer, nil
}

// Close shuts down the synchronous Kafka producer, handling potential errors.
func (p *ProducerKafka) Close() {
	if err := p.saramaProducer.Close(); err != nil {
		logrus.Errorf("Error closing producer: %v", err)
	}
}

// SandMsgToKafka SendMsgToKafka sends a message to a specified Kafka topic using the synchronous producer.
// It takes the topic name, request ID, and data as parameters and returns an error if encountered.
func (p *ProducerKafka) SandMsgToKafka(ctx context.Context, topicName, requestID string, data []byte) error {
	// Формируем с данными полученными из запроса
	responseMsg := &sarama.ProducerMessage{
		Topic: "Handler",
		Key:   sarama.StringEncoder(requestID),
		Value: sarama.StringEncoder(data),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("topic_name"),
				Value: []byte(topicName),
			},
		},
	}
	select {
	case <-ctx.Done():
		logrus.Info("Context canceled, not sending message to Kafka")
		return ctx.Err()
	default:
		// Отправляем сообщение в Kafka
		_, _, err := p.saramaProducer.SendMessage(responseMsg)
		if err != nil {
			logrus.Errorf("Failed to send message to Kafka: %v", err)
			return err
		}
		logrus.Infof("Sended message: %v", responseMsg)
		return nil
	}
}
