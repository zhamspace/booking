package expiry

import (
	"context"
	"log/slog"
	"time"

	bookingConstants "github.com/zhamspace/booking/internal/domain/booking/constants"
	bookingModel "github.com/zhamspace/booking/internal/domain/booking/model"
	commonModel "github.com/zhamspace/booking/internal/domain/common/model"
)

type BookingServiceI interface {
	List(ctx context.Context, pars *bookingModel.ListReq) ([]*bookingModel.Main, int64, error)
	Update(ctx context.Context, pars *bookingModel.GetReq, obj *bookingModel.Edit) error
}

type Worker struct {
	bookingService BookingServiceI
	interval       time.Duration
}

func New(bookingService BookingServiceI, interval time.Duration) *Worker {
	return &Worker{bookingService: bookingService, interval: interval}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.tick(ctx)
		}
	}
}

func (w *Worker) tick(ctx context.Context) {
	now := time.Now().UTC()
	statusCreated := bookingConstants.StatusCreated
	items, _, err := w.bookingService.List(ctx, &bookingModel.ListReq{
		ListParams: commonModel.ListParams{PageSize: 200},
		Status:     &statusCreated,
	})
	if err != nil {
		slog.Error("expiry worker: list bookings", "error", err)
		return
	}

	statusCancelled := bookingConstants.StatusCancelled
	cancelReason := "payment_timeout"

	for _, b := range items {
		if b.HoldExpiresAt == nil || !b.HoldExpiresAt.Before(now) {
			continue
		}
		err := w.bookingService.Update(ctx, &bookingModel.GetReq{Id: b.Id}, &bookingModel.Edit{
			Status:       &statusCancelled,
			CancelReason: &cancelReason,
			CancelledAt:  &now,
			UpdatedAt:    &now,
		})
		if err != nil {
			slog.Error("expiry worker: cancel booking", "booking_id", b.Id, "error", err)
		} else {
			slog.Info("expiry worker: cancelled stale booking", "booking_id", b.Id)
		}
	}
}
