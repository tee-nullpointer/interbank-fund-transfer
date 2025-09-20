package service

import (
	"context"
	"encoding/json"
	"iso8583-gateway/internal/config"
	"iso8583-gateway/internal/domain"
	"iso8583-gateway/pkg/util"
	"sync"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type InboundService struct {
	ctx               context.Context
	inboundChan       chan *domain.ISO8583Message
	applicationConfig *config.ApplicationConfig
	producer          sarama.SyncProducer
	serviceID         string
	wg                sync.WaitGroup
}

func NewInboundService(ctx context.Context, inboundChan chan *domain.ISO8583Message, applicationConfig *config.ApplicationConfig, producer sarama.SyncProducer) *InboundService {
	var serviceID string
	id, err := uuid.NewUUID()
	if err != nil {
		serviceID = util.RandomString(36)
	} else {
		serviceID = id.String()
	}
	return &InboundService{
		ctx:               ctx,
		inboundChan:       inboundChan,
		applicationConfig: applicationConfig,
		producer:          producer,
		serviceID:         serviceID,
	}
}

func (service *InboundService) Start() {
	for range service.applicationConfig.WorkerPerConnection {
		service.wg.Add(1)
		go func() {
			defer service.wg.Done()
			service.processInbound()
		}()
	}
	service.wg.Wait()
}

func (service *InboundService) processInbound() {
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
			de63 := v.Fields[63]
			if de63 == "" {
				zap.L().Warn("Ignore message with empty DE63", zap.Any("fields", v.Fields))
				continue
			}
			bytes, err := json.Marshal(v)
			if err != nil {
				zap.L().Error("Failed to marshal ISO8583Message to JSON", zap.Error(err), zap.Any("message", v), zap.String("de63", de63))
				continue
			}
			msg := &sarama.ProducerMessage{
				Topic: service.applicationConfig.InboundRequestTopic,
				Value: sarama.ByteEncoder(bytes),
				Headers: []sarama.RecordHeader{
					{Key: []byte("service_id"), Value: []byte(service.serviceID)},
					{Key: []byte("de63"), Value: []byte(de63)},
				},
			}
			partition, offset, err := service.producer.SendMessage(msg)
			if err != nil {
				zap.L().Error("Failed to send message to Kafka", zap.Error(err), zap.Any("message", v), zap.String("de63", de63))
				continue
			}
			zap.L().Info("Successfully sent message to Kafka", zap.Any("message", v), zap.String("de63", de63), zap.Int64("offset", offset), zap.Int32("partition", partition))
		}
	}
}
