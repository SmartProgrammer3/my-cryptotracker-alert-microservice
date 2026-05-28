package sniper

import (
	"cryptotracker/alert/internal/alert"
	"cryptotracker/alert/internal/db"
	"log"
)

// Add adiciona um novo alerta à memória.
func (s *Sniper) Add(a alert.Alert) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.alerts[a.Symbol] = append(s.alerts[a.Symbol], a)
}

// Check avalia o preço contra os alertas activos do símbolo.
// Retorna true se o símbolo não tiver mais alertas PENDING.
func (s *Sniper) Check(symbol string, price float64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	alerts := s.alerts[symbol]
	remaining := alerts[:0]

	for _, a := range alerts {
		if price >= a.TargetPrice {
			s.done(a)
		} else {
			remaining = append(remaining, a)
		}
	}

	s.alerts[symbol] = remaining
	return len(remaining) == 0
}

// done marca o alerta como TRIGGERED na DB e remove da memória.
func (s *Sniper) done(a alert.Alert) {
	if err := db.UpdateAlertStatus(s.db, a.ID, alert.StatusTriggered); err != nil {
		log.Printf("[sniper] erro ao marcar alerta %d como TRIGGERED: %v", a.ID, err)
		return
	}
	log.Printf("[sniper] alerta %d disparado — %s atingiu %.2f", a.ID, a.Symbol, a.TargetPrice)
}

// Load carrega alertas PENDING da DB para memória no arranque do serviço.
func (s *Sniper) Load(alerts []alert.Alert) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, a := range alerts {
		s.alerts[a.Symbol] = append(s.alerts[a.Symbol], a)
	}

	log.Printf("[sniper] %d alertas PENDING carregados em memória", len(alerts))
}

// Remove elimina um alerta da memória (cancelado via gRPC).
// Retorna o símbolo e true se ficou sem alertas PENDING.
func (s *Sniper) Remove(id int) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for symbol, alerts := range s.alerts {
		remaining := alerts[:0]
		for _, a := range alerts {
			if a.ID != id {
				remaining = append(remaining, a)
			}
		}
		s.alerts[symbol] = remaining

		if len(remaining) == 0 {
			return symbol, true
		}
	}

	return "", false
}