package producers

import (
	"github.com/IBM/sarama"
	"log"
)

func ProduceMessage(brokerAddress, topicName string, data []byte) error {
	// Создание продюсера Kafka
	producer, err := sarama.NewSyncProducer([]string{brokerAddress}, nil)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Формируем ответное сообщение
	resp := &sarama.ProducerMessage{
		Topic: "Handler",
		Key:   sarama.StringEncoder("Key1"),
		Value: sarama.StringEncoder(data),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("topic_name"),
				Value: []byte(topicName),
			},
		},
	}
	// Отправляем сообщение в Kafka
	_, _, err = producer.SendMessage(resp)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
	}
	return nil
}
