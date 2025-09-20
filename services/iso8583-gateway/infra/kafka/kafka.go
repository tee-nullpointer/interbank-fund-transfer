package kafka

import (
	"iso8583-gateway/internal/config"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

var producer sarama.SyncProducer

func InitKafka(cfg *config.KafkaConfig) sarama.SyncProducer {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.RequiredAcks = sarama.WaitForAll
	saramaCfg.Producer.Retry.Max = cfg.Retry
	saramaCfg.Producer.Return.Successes = true
	saramaCfg.Producer.Return.Errors = true
	saramaCfg.Producer.Timeout = cfg.Timeout

	syncProducer, err := sarama.NewSyncProducer(cfg.Brokers, saramaCfg)
	if err != nil {
		panic(err)
	}
	producer = syncProducer
	return producer
}

func Close() {
	if producer != nil {
		if err := producer.Close(); err != nil {
			zap.L().Error("Fail to close kafka producer", zap.Error(err))
		}
	}
}
