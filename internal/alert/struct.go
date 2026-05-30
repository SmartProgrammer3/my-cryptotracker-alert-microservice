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

type Direction string

const (
    DirectionAbove Direction = "ABOVE"
    DirectionBelow Direction = "BELOW"
)

type Alert struct {
	ID          int
	Symbol      string
	TargetPrice float64
	Direction   Direction
	Status      Status
	CreatedAt   time.Time
	TriggeredAt *time.Time
}
