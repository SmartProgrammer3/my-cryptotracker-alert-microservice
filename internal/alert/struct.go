package alert

import (
	"time"
)

type Status string

const (
    StatusPending   Status = "PENDING"
    StatusTriggered Status = "TRIGGERED"
    StatusCancelled Status = "CANCELLED"
)

type Alert struct {
	ID          int
	Symbol      string
	TargetPrice float64
	Status      Status
	CreatedAt   time.Time
	TriggeredAt *time.Time
}
