package alert

import "time"

func NewAlert(id int, symbol string, targetPrice float64, createdAt time.Time) Alert {
	return Alert{
		ID:          id,
		Symbol:      symbol,
		TargetPrice: targetPrice,
		Status:      StatusPending,
		CreatedAt:   createdAt,
	}
}