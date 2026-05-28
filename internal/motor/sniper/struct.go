package sniper

import (
	"database/sql"
	"sync"

	"cryptotracker/alert/internal/alert"
)

type Sniper struct {
	mu     sync.Mutex
	alerts map[string][]alert.Alert
	db     *sql.DB
}
