package dto

import (
	"strings"
	"time"

	bookingModel "github.com/zhamspace/booking/internal/domain/booking/model"
	"github.com/zhamspace/booking/pkg/proto/booking_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func EncodeBookingMain(v *bookingModel.Main, _ int) *booking_v1.BookingMain {
	res := &booking_v1.BookingMain{
		Id:              v.Id,
		VenueId:         v.VenueId,
		ResourceId:      v.ResourceId,
		SessionId:       v.SessionId,
		UserId:          v.UserId,
		Status:          v.Status,
		PaymentStatus:   v.PaymentStatus,
		PriceTotal:      v.PriceTotal,
		Currency:        v.Currency,
		Timezone:        v.Timezone,
		CreatedAt:       timestamppb.New(v.CreatedAt),
		UpdatedAt:       timestamppb.New(v.UpdatedAt),
		PaymentIntentId: v.PaymentIntentId,
		CancelReason:    v.CancelReason,
	}

	if v.StartAt != nil {
		res.StartAt = timestamppb.New(*v.StartAt)
	}
	if v.EndAt != nil {
		res.EndAt = timestamppb.New(*v.EndAt)
	}
	if v.HoldExpiresAt != nil {
		res.HoldExpiresAt = timestamppb.New(*v.HoldExpiresAt)
	}
	if v.CancelledAt != nil {
		res.CancelledAt = timestamppb.New(*v.CancelledAt)
	}
	if v.ConfirmedAt != nil {
		res.ConfirmedAt = timestamppb.New(*v.ConfirmedAt)
	}
	if v.CompletedAt != nil {
		res.CompletedAt = timestamppb.New(*v.CompletedAt)
	}
	return res
}

func DecodeOptionalString(v *string) *string {
	if v == nil {
		return nil
	}

	s := strings.TrimSpace(*v)
	if s == "" {
		return nil
	}

	return &s
}

func DecodeTimestampPtr(v *timestamppb.Timestamp) *time.Time {
	if v == nil {
		return nil
	}

	return new(v.AsTime())
}
