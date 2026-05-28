package scout

import (
	"context"
	"sync"

	"cryptotracker/alert/internal/motor/sniper"
)

type Scout struct {
	mu           sync.Mutex
	streams      map[string]context.CancelFunc
	sniper       *sniper.Sniper
	binanceWSURL string
}
