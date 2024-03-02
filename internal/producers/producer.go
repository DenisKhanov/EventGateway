package producers

//
//import (
//	"context"
//	"github.com/segmentio/kafka-go"
//	"github.com/sirupsen/logrus"
//)
//
//func ProduceMessage(ctx context.Context, brokerAddress, topicName string, data []byte) error {
//	// Создаем писателя для записи в топик
//	w := kafka.Writer{
//		Addr:     kafka.TCP(brokerAddress),
//		Balancer: &kafka.LeastBytes{},
//		Topic:    "Handler",
//	}
//
//	// Создаем сообщение для отправки
//	message := kafka.Message{
//		Key:   []byte("key1"), // Ключ сообщения (может быть nil)
//		Value: data,           // Значение сообщения
//		Headers: []kafka.Header{
//			{Key: "topic_name", Value: []byte(topicName)},
//		},
//	}
//
//	// Публикуем сообщение в топик
//	err := w.WriteMessages(ctx, message)
//	if err != nil {
//		logrus.Fatal("Error producing message:", err)
//		return err
//	}
//
//	// Закрываем писатель после использования
//	w.Close()
//	return nil
//}
