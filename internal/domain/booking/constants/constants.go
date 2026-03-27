package constants

import "time"

const (
	StatusCreated   string = "created"
	StatusConfirmed string = "confirmed"
	StatusCancelled string = "cancelled"
	StatusCompleted string = "completed"
)

const (
	PaymentStatusPending string = "pending"
	PaymentStatusPaid    string = "paid"
	PaymentStatusFailed  string = "failed"
)

const (
	DefaultCurrency     string        = "KZT"
	DefaultHoldDuration time.Duration = 15 * time.Minute
	DefaultCapacity     int64         = 1
)
