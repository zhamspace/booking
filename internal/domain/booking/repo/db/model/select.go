package model

import (
	"time"

	bookingModel "github.com/zhamspace/booking/internal/domain/booking/model"
)

type Select struct {
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

func (m *Select) ListColumnMap() map[string]any {
	return map[string]any{
		"id":                &m.Id,
		"venue_id":          &m.VenueId,
		"resource_id":       &m.ResourceId,
		"session_id":        &m.SessionId,
		"user_id":           &m.UserId,
		"created_at":        &m.CreatedAt,
		"updated_at":        &m.UpdatedAt,
		"status":            &m.Status,
		"payment_status":    &m.PaymentStatus,
		"price_total":       &m.PriceTotal,
		"currency":          &m.Currency,
		"start_at":          &m.StartAt,
		"end_at":            &m.EndAt,
		"timezone":          &m.Timezone,
		"hold_expires_at":   &m.HoldExpiresAt,
		"payment_intent_id": &m.PaymentIntentId,
		"cancel_reason":     &m.CancelReason,
		"cancelled_at":      &m.CancelledAt,
		"confirmed_at":      &m.ConfirmedAt,
		"completed_at":      &m.CompletedAt,
	}
}

func (m *Select) PKColumnMap() map[string]any {
	return map[string]any{
		"id": &m.Id,
	}
}

func (m *Select) DefaultSortColumns() []string {
	return []string{
		"created_at DESC",
	}
}

func EncodeMain(src *Select, _ int) *bookingModel.Main {
	if src == nil {
		return nil
	}

	return &bookingModel.Main{
		Id:              src.Id,
		VenueId:         src.VenueId,
		ResourceId:      src.ResourceId,
		SessionId:       src.SessionId,
		UserId:          src.UserId,
		CreatedAt:       src.CreatedAt,
		UpdatedAt:       src.UpdatedAt,
		Status:          src.Status,
		PaymentStatus:   src.PaymentStatus,
		PriceTotal:      src.PriceTotal,
		Currency:        src.Currency,
		StartAt:         src.StartAt,
		EndAt:           src.EndAt,
		Timezone:        src.Timezone,
		HoldExpiresAt:   src.HoldExpiresAt,
		PaymentIntentId: src.PaymentIntentId,
		CancelReason:    src.CancelReason,
		CancelledAt:     src.CancelledAt,
		ConfirmedAt:     src.ConfirmedAt,
		CompletedAt:     src.CompletedAt,
	}
}
