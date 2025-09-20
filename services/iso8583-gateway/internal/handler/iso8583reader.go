package handler

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"iso8583-gateway/internal/domain"
	"iso8583-gateway/pkg/util"
	"net"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type ISO8583Reader struct {
	conn        net.Conn
	ctx         context.Context
	inboundChan chan *domain.ISO8583Message
}

func NewISO8583Reader(conn net.Conn, ctx context.Context, inboundChan chan *domain.ISO8583Message) *ISO8583Reader {
	return &ISO8583Reader{
		conn:        conn,
		ctx:         ctx,
		inboundChan: inboundChan,
	}
}

func (reader *ISO8583Reader) Read() {
	defer closeConnection(reader.conn)
	r := bufio.NewReader(reader.conn)
	remoteAddress := reader.conn.RemoteAddr().String()
	for {
		if reader.isShuttingDown() {
			zap.L().Info("Server shutting down, closing connection", zap.String("remote_addr", remoteAddress))
			return
		}
		msgLen, err := reader.readMessageLength(r, remoteAddress)
		if err != nil {
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				continue
			}
			return
		}
		msg, err := reader.readMessage(r, remoteAddress, msgLen)
		if err != nil {
			return
		}
		zap.L().Info("ISO8583 message parsed", zap.String("remote_addr", remoteAddress), zap.String("mti", msg.MTI), zap.Any("fields", msg.Fields))
		reader.inboundChan <- msg
	}
}

func (reader *ISO8583Reader) isShuttingDown() bool {
	select {
	case <-reader.ctx.Done():
		return true
	default:
		return false
	}
}

func (reader *ISO8583Reader) readMessageLength(r *bufio.Reader, remoteAddress string) (int, error) {
	lenBuf := make([]byte, 4)
	err := reader.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		zap.L().Error("Error setting read deadline", zap.String("remote_addr", remoteAddress), zap.Error(err))
		return 0, err
	}
	_, err = io.ReadFull(r, lenBuf)
	if err != nil {
		var ne net.Error
		if errors.As(err, &ne) && ne.Timeout() {
			return 0, err
		}
		zap.L().Info("Client disconnected", zap.String("remote_addr", remoteAddress))
		return 0, err
	}
	msgLen, err := strconv.Atoi(string(lenBuf))
	if err != nil || msgLen < 0 || msgLen > 2048 {
		zap.L().Warn("Invalid header length, closing connection", zap.String("remote_addr", remoteAddress), zap.Error(err))
		return 0, err
	}
	return msgLen, nil
}

func (reader *ISO8583Reader) readMessage(r *bufio.Reader, remoteAddress string, msgLen int) (*domain.ISO8583Message, error) {
	msgBuf := make([]byte, msgLen)
	_, err := io.ReadFull(r, msgBuf)
	if err != nil {
		zap.L().Info("Invalid message body, closing connection", zap.String("remote_addr", remoteAddress), zap.Error(err))
		return nil, err
	}
	zap.L().Info("Raw message received", zap.String("remote_addr", remoteAddress), zap.String("raw_message", fmt.Sprintf("% X", msgBuf)))
	msg, err := util.ParseISO8583(msgBuf)
	if err != nil {
		zap.L().Error("Fail to parse ISO8583 message", zap.String("remote_addr", remoteAddress), zap.Error(err))
		return nil, err
	}
	return msg, nil
}

func closeConnection(conn net.Conn) {
	if err := conn.Close(); err != nil {
		zap.L().Error("Error closing connection", zap.Error(err))
	}
}
