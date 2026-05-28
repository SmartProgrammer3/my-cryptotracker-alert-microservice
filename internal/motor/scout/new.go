package scout

import (
	"context"

	"cryptotracker/alert/internal/motor/sniper"
)

func New(sniper *sniper.Sniper, binanceWSURL string) *Scout {
	return &Scout{
		streams:      make(map[string]context.CancelFunc),
		sniper:       sniper,
		binanceWSURL: binanceWSURL,
	}
}