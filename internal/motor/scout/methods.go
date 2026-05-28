package scout

import (
	"context"
	"log"
)

// Subscribe abre um stream WS para o símbolo, se ainda não existir.
func (s *Scout) Subscribe(symbol string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.streams[symbol]; exists {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.streams[symbol] = cancel

	log.Printf("[scout] a subscrever stream para %s", symbol)
	go s.runStream(ctx, symbol)
}

// Unsubscribe fecha o stream WS do símbolo se não houver mais alertas PENDING para o símbolo.
func (s *Scout) Unsubscribe(symbol string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if cancel, exists := s.streams[symbol]; exists {
        cancel()
        delete(s.streams, symbol)
        log.Printf("[scout] stream %s fechado", symbol)
    }
}