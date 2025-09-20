package main

import (
	"iso8583-gateway/infra/kafka"
	"iso8583-gateway/internal/config"
	"iso8583-gateway/internal/server"
	"iso8583-gateway/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Init()
	err := logger.InitLogger(cfg.Logger.Level, cfg.Logger.Format)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	producer := kafka.InitKafka(cfg.Kafka)
	defer kafka.Close()
	srv := server.NewServer(cfg.Server.Host, cfg.Server.Port, cfg.Application, producer)
	go srv.Start()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	zap.L().Info("Server is shutting down...")
	srv.Shutdown()
	zap.L().Info("Server stopped")
}
