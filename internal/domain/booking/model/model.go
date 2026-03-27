package model

import (
	"time"

	commonModel "github.com/zhamspace/booking/internal/domain/common/model"
)

type Main struct {
	Id         string
	VenueId    string
	ResourceId string
	SessionId  *string
	UserId     string
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Status          string
	PaymentStatus   string
	PriceTotal      int64
	Currency        string
	StartAt         *time.Time
	EndAt           *time.Time
	Timezone        string
	HoldExpiresAt   *time.Time
	PaymentIntentId *string
	CancelReason    *string
	CancelledAt     *time.Time
	ConfirmedAt     *time.Time
	CompletedAt     *time.Time
}

type Edit struct {
	Id         *string
	VenueId    *string
	ResourceId *string
	SessionId  *string
	UserId     *string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time

	Status          *string
	PaymentStatus   *string
	PriceTotal      *int64
	Currency        *string
	StartAt         *time.Time
	EndAt           *time.Time
	Timezone        *string
	HoldExpiresAt   *time.Time
	PaymentIntentId *string
	CancelReason    *string
	CancelledAt     *time.Time
	ConfirmedAt     *time.Time
	CompletedAt     *time.Time
}

type ListReq struct {
	commonModel.ListParams

	Ids              *[]string
	VenueId          *string
	ResourceId       *string
	UserId           *string
	SessionId        *string
	Status           *string
	PaymentStatus    *string
	Currency         *string
	From             *time.Time
	To               *time.Time
	IncludeCancelled bool
	ActiveAt         *time.Time
}

type GetReq struct {
	Id string
}
