package server

import (
	"context"
	"fmt"
	"iso8583-gateway/internal/config"
	"iso8583-gateway/internal/domain"
	"iso8583-gateway/internal/handler"
	"iso8583-gateway/internal/service"
	"net"
	"sync"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Server struct {
	address    string
	listener   net.Listener
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
	cfg        *config.ApplicationConfig
	producer   sarama.SyncProducer
}

func NewServer(host string, port string, cfg *config.ApplicationConfig, producer sarama.SyncProducer) *Server {
	address := fmt.Sprintf("%s:%s", host, port)
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		address:    address,
		ctx:        ctx,
		cancelFunc: cancel,
		cfg:        cfg,
		producer:   producer,
	}
}

func (server *Server) Start() {
	zap.L().Info("Starting server", zap.String("address", server.address))
	ln, err := net.Listen("tcp", server.address)
	if err != nil {
		zap.L().Fatal("Failed to start server", zap.Error(err))
	}
	server.listener = ln
	zap.L().Info("Server started", zap.String("address", server.address))
	go server.acceptConnection()
}

func (server *Server) acceptConnection() {
	zap.L().Info("Accepting connection", zap.String("address", server.address))
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			select {
			case <-server.ctx.Done():
				return
			default:
				zap.L().Error("Failed to accept connection", zap.Error(err))
				continue
			}
		}
		zap.L().Info("New connection accepted", zap.String("remote_addr", conn.RemoteAddr().String()))
		inboundChan := make(chan *domain.ISO8583Message, 200)
		server.wg.Add(2)
		reader := handler.NewISO8583Reader(conn, server.ctx, inboundChan)
		inboundService := service.NewInboundService(server.ctx, inboundChan, server.cfg, server.producer)
		go func() {
			defer server.wg.Done()
			reader.Read()
		}()
		go func() {
			defer server.wg.Done()
			inboundService.Start()
		}()
	}
}

func (server *Server) Shutdown() {
	zap.L().Info("Shutting down server")
	server.cancelFunc()
	if server.listener != nil {
		if err := server.listener.Close(); err != nil {
			zap.L().Error("Failed to close server listener", zap.Error(err))
		}
	}
	server.wg.Wait()
	zap.L().Info("Server shutdown completed")
}
