package model

import (
	"time"

	bookingModel "github.com/zhamspace/booking/internal/domain/booking/model"
)

type Upsert struct {
	PKId string

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

func (m *Upsert) CreateColumnMap() map[string]any {
	result := map[string]any{}

	if m.Id != nil {
		result["id"] = *m.Id
	}
	if m.VenueId != nil {
		result["venue_id"] = *m.VenueId
	}
	if m.ResourceId != nil {
		result["resource_id"] = *m.ResourceId
	}
	if m.SessionId != nil {
		result["session_id"] = *m.SessionId
	}
	if m.UserId != nil {
		result["user_id"] = *m.UserId
	}
	if m.CreatedAt != nil {
		result["created_at"] = *m.CreatedAt
	}
	if m.UpdatedAt != nil {
		result["updated_at"] = *m.UpdatedAt
	}
	if m.Status != nil {
		result["status"] = *m.Status
	}
	if m.PaymentStatus != nil {
		result["payment_status"] = *m.PaymentStatus
	}
	if m.PriceTotal != nil {
		result["price_total"] = *m.PriceTotal
	}
	if m.Currency != nil {
		result["currency"] = *m.Currency
	}
	if m.StartAt != nil {
		result["start_at"] = *m.StartAt
	}
	if m.EndAt != nil {
		result["end_at"] = *m.EndAt
	}
	if m.Timezone != nil {
		result["timezone"] = *m.Timezone
	}
	if m.HoldExpiresAt != nil {
		result["hold_expires_at"] = *m.HoldExpiresAt
	}
	if m.PaymentIntentId != nil {
		result["payment_intent_id"] = *m.PaymentIntentId
	}
	if m.CancelReason != nil {
		result["cancel_reason"] = *m.CancelReason
	}
	if m.CancelledAt != nil {
		result["cancelled_at"] = *m.CancelledAt
	}
	if m.ConfirmedAt != nil {
		result["confirmed_at"] = *m.ConfirmedAt
	}
	if m.CompletedAt != nil {
		result["completed_at"] = *m.CompletedAt
	}

	return result
}

func (m *Upsert) ReturningColumnMap() map[string]any {
	return map[string]any{
		"id": &m.PKId,
	}
}

func (m *Upsert) UpdateColumnMap() map[string]any {
	result := map[string]any{}

	if m.VenueId != nil {
		result["venue_id"] = *m.VenueId
	}
	if m.ResourceId != nil {
		result["resource_id"] = *m.ResourceId
	}
	if m.SessionId != nil {
		result["session_id"] = *m.SessionId
	}
	if m.UserId != nil {
		result["user_id"] = *m.UserId
	}
	if m.UpdatedAt != nil {
		result["updated_at"] = *m.UpdatedAt
	}
	if m.Status != nil {
		result["status"] = *m.Status
	}
	if m.PaymentStatus != nil {
		result["payment_status"] = *m.PaymentStatus
	}
	if m.PriceTotal != nil {
		result["price_total"] = *m.PriceTotal
	}
	if m.Currency != nil {
		result["currency"] = *m.Currency
	}
	if m.StartAt != nil {
		result["start_at"] = *m.StartAt
	}
	if m.EndAt != nil {
		result["end_at"] = *m.EndAt
	}
	if m.Timezone != nil {
		result["timezone"] = *m.Timezone
	}
	if m.HoldExpiresAt != nil {
		result["hold_expires_at"] = *m.HoldExpiresAt
	}
	if m.PaymentIntentId != nil {
		result["payment_intent_id"] = *m.PaymentIntentId
	}
	if m.CancelReason != nil {
		result["cancel_reason"] = *m.CancelReason
	}
	if m.CancelledAt != nil {
		result["cancelled_at"] = *m.CancelledAt
	}
	if m.ConfirmedAt != nil {
		result["confirmed_at"] = *m.ConfirmedAt
	}
	if m.CompletedAt != nil {
		result["completed_at"] = *m.CompletedAt
	}

	return result
}

func (m *Upsert) PKColumnMap() map[string]any {
	return map[string]any{
		"id": m.PKId,
	}
}

func DecodeEdit(src *bookingModel.Edit) *Upsert {
	if src == nil {
		return &Upsert{}
	}

	return &Upsert{
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
