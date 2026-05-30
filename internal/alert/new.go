package alert

import "time"

func NewAlert(id int, symbol string, targetPrice float64, direction Direction, createdAt time.Time) Alert {
	return Alert{
		ID:          id,
		Symbol:      symbol,
		TargetPrice: targetPrice,
		Direction:   direction,
		Status:      StatusPending,
		CreatedAt:   createdAt,
	}
}