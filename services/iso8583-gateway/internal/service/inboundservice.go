package service

import (
	"context"
	"encoding/json"
	"iso8583-gateway/internal/config"
	"iso8583-gateway/internal/domain"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type InboundService struct {
	ctx               context.Context
	inboundChan       chan *domain.ISO8583Message
	applicationConfig *config.ApplicationConfig
	producer          sarama.SyncProducer
}

func NewInboundService(ctx context.Context, inboundChan chan *domain.ISO8583Message, applicationConfig *config.ApplicationConfig, producer sarama.SyncProducer) *InboundService {
	return &InboundService{
		ctx:               ctx,
		inboundChan:       inboundChan,
		applicationConfig: applicationConfig,
		producer:          producer,
	}
}

func (service *InboundService) ProcessInbound() {
	for {
		select {
		case <-service.ctx.Done():
			zap.L().Info("Context is done, stopping service")
			return
		case v, ok := <-service.inboundChan:
			if !ok {
				zap.L().Info("Inbound channel is closed, stopping service")
				return
			}
			go service.processInbound(v)
		}
	}
}

func (service *InboundService) processInbound(v *domain.ISO8583Message) {
	f63 := v.Fields[63]
	if f63 == "" {
		zap.L().Warn("Ignore message with empty F63", zap.Any("fields", v.Fields))
		return
	}
	bytes, err := json.Marshal(v)
	if err != nil {
		zap.L().Error("Failed to marshal ISO8583Message to JSON", zap.Error(err), zap.Any("message", v), zap.String("f63", f63))
		return
	}
	msg := &sarama.ProducerMessage{
		Topic: service.applicationConfig.InboundRequestTopic,
		Value: sarama.ByteEncoder(bytes),
		Headers: []sarama.RecordHeader{
			{Key: []byte("service_id"), Value: []byte(service.applicationConfig.ServiceID)},
			{Key: []byte("trace_id"), Value: []byte(f63)},
		},
	}
	partition, offset, err := service.producer.SendMessage(msg)
	if err != nil {
		zap.L().Error("Failed to send message to Kafka", zap.Error(err), zap.Any("message", v), zap.String("f63", f63))
		return
	}
	zap.L().Info("Successfully sent message to Kafka", zap.Any("message", v), zap.String("f63", f63), zap.Int64("offset", offset), zap.Int32("partition", partition))
}
