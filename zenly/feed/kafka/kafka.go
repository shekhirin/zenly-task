package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/shekhirin/zenly-task/zenly/feed"
	"github.com/shekhirin/zenly-task/zenly/pb"
	log "github.com/sirupsen/logrus"
)

type kafkaFeed struct {
	producer sarama.AsyncProducer
	topic    string
}

func New(producer sarama.AsyncProducer, topic string) feed.Feed {
	go func() {
		for err := range producer.Errors() {
			log.WithError(err).Info("failed to write to kafka producer")
		}
	}()

	return &kafkaFeed{
		producer: producer,
		topic:    topic,
	}
}

func (feed kafkaFeed) Publish(message *pb.FeedMessage) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	feed.producer.Input() <- &sarama.ProducerMessage{
		Topic: feed.topic,
		Value: sarama.ByteEncoder(data),
	}

	return nil
}
