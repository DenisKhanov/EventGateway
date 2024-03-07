package consumers

import (
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"log"
)

type ConsumerKafka struct {
	brokerAddress  string
	saramaConsumer sarama.Consumer
}

// NewConsumerKafka creates a new instance of ConsumerKafka, which represents a Kafka consumer.
// It takes the broker address as a parameter and returns a pointer to ConsumerKafka and any error encountered during creation.
func NewConsumerKafka(brokerAddress string) (*ConsumerKafka, error) {
	saramaConsumer, err := createSaramaConsumer(brokerAddress)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &ConsumerKafka{
		brokerAddress:  brokerAddress,
		saramaConsumer: saramaConsumer,
	}, nil
}

// createSaramaConsumer creates a Sarama Kafka consumer.
// It takes the broker address as a parameter and returns the created consumer.
func createSaramaConsumer(brokerAddress string) (sarama.Consumer, error) {
	// Создание консьюмера Kafka
	saramaConsumer, err := sarama.NewConsumer([]string{brokerAddress}, nil)
	if err != nil {
		logrus.Errorf("Failed to create consumer: %v", err)
		return nil, err
	}
	return saramaConsumer, nil
}

// Close closes the Sarama Kafka consumer.
// It ensures that the consumer is properly closed, releasing any allocated resources.
func (c *ConsumerKafka) Close() {
	// Закрытие консьюмера Kafka
	if err := c.saramaConsumer.Close(); err != nil {
		log.Printf("Error closing consumer: %v", err)
	}
}

// ListenKafkaMessages subscribes to a specified Kafka topic and listens for incoming messages.
// It takes the topic name and a channel for receiving messages as parameters.
// The method continuously processes incoming messages and sends them to the provided channel.
func (c *ConsumerKafka) ListenKafkaMessages(topicName string, message chan<- *sarama.ConsumerMessage) error {
	// Подписка на определенную партицию топика в Kafka
	partConsumer, err := c.saramaConsumer.ConsumePartition(topicName, 0, sarama.OffsetNewest)
	if err != nil {
		logrus.Errorf("Failed to consume partition: %v, topic_name : %s", err, topicName)
		return err
	}
	defer func() {
		if err = partConsumer.Close(); err != nil {
			logrus.Errorf("Error closing partition consumer: %v", err)
		}
		logrus.Info("ConsumerKafka closed")
	}()
	for {
		select {
		// получение входящего сообщения из Kafka и отправка в service
		case msg, ok := <-partConsumer.Messages():
			if !ok {
				logrus.Infof("Channel closed, exiting")
				return nil
			}
			logrus.Infof("Received message: %+v\n", string(msg.Value))
			message <- msg
		}
	}
}
