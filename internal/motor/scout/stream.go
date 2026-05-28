package scout

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type binancePriceEvent struct {
    Price string `json:"p"`
}

func (s *Scout) runStream(ctx context.Context, symbol string) {
	url := fmt.Sprintf("%s/%s@aggTrade", s.binanceWSURL, symbol)
	backoff := time.Second

	for {
		select {
		case <-ctx.Done():
			log.Printf("[scout] stream %s terminado", symbol)
			return
		default:
		}

		conn, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
		if err != nil {
			log.Printf("[scout] erro ao conectar %s: %v — retry em %s", symbol, err, backoff)
			time.Sleep(backoff)
			backoff = min(backoff*2, 30*time.Second)
			continue
		}

		backoff = time.Second
		log.Printf("[scout] stream %s conectado", symbol)

		s.readStream(ctx, conn, symbol)
	}
}

func (s *Scout) readStream(ctx context.Context, conn *websocket.Conn, symbol string) {
	defer conn.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[scout] stream %s desconectado: %v", symbol, err)
			return
		}

		var event binancePriceEvent
		if err := json.Unmarshal(msg, &event); err != nil {
			continue
		}

		price, err := strconv.ParseFloat(event.Price, 64)
		if err != nil {
			continue
		}

		if s.sniper.Check(symbol, price) {
			s.Unsubscribe(symbol)
		}
	}
}