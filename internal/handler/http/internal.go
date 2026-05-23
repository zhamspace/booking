package http

import (
	"net/http"
	"strings"
	"time"

	bookingConstants "github.com/zhamspace/booking/internal/domain/booking/constants"
	bookingModel "github.com/zhamspace/booking/internal/domain/booking/model"
	usecaseBookingP "github.com/zhamspace/booking/internal/usecase/booking"
)

type Internal struct {
	bookingUsecase *usecaseBookingP.Usecase
}

func NewInternal(bookingUsecase *usecaseBookingP.Usecase) *Internal {
	return &Internal{bookingUsecase: bookingUsecase}
}

// GET /internal/bookings/{id}
func (h *Internal) GetBooking(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	id := strings.TrimSpace(pathParams["id"])
	if id == "" {
		errorResponse(w, http.StatusBadRequest, nil)
		return
	}

	item, err := h.bookingUsecase.Get(r.Context(), &bookingModel.GetReq{Id: id})
	if err != nil {
		errorResponse(w, 0, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":             item.Id,
		"user_id":        item.UserId,
		"venue_id":       item.VenueId,
		"resource_id":    item.ResourceId,
		"status":         item.Status,
		"payment_status": item.PaymentStatus,
		"price_total":    item.PriceTotal,
		"currency":       item.Currency,
	})
}

// POST /internal/bookings/{id}/mark-paid
// Body: {"payment_intent_id": "pi_mock_..."}
func (h *Internal) MarkPaid(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	id := strings.TrimSpace(pathParams["id"])
	if id == "" {
		errorResponse(w, http.StatusBadRequest, nil)
		return
	}

	var body struct {
		PaymentIntentId string `json:"payment_intent_id"`
	}
	if err := decodeJSON(r, &body); err != nil {
		errorResponse(w, 0, err)
		return
	}

	intentId := body.PaymentIntentId
	item, err := h.bookingUsecase.Confirm(r.Context(), &bookingModel.ConfirmReq{
		Id:              id,
		PaymentIntentId: &intentId,
	})
	if err != nil {
		errorResponse(w, 0, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": item.Id, "status": item.Status})
}

// POST /internal/bookings/{id}/mark-failed
func (h *Internal) MarkFailed(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	id := strings.TrimSpace(pathParams["id"])
	if id == "" {
		errorResponse(w, http.StatusBadRequest, nil)
		return
	}

	now := time.Now().UTC()
	paymentStatusFailed := bookingConstants.PaymentStatusFailed
	if err := h.bookingUsecase.UpdatePaymentStatus(r.Context(), id, &paymentStatusFailed, &now); err != nil {
		errorResponse(w, 0, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /internal/bookings/stats?venue_ids=id1,id2&from=RFC3339&to=RFC3339
// Internal — same aggregates as public /bookings/stats, for trusted services (e.g. payment).
func (h *Internal) Stats(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	raw := strings.TrimSpace(r.URL.Query().Get("venue_ids"))
	if raw == "" {
		errorResponse(w, http.StatusBadRequest, nil)
		return
	}
	var venueIDs []string
	for _, p := range strings.Split(raw, ",") {
		id := strings.TrimSpace(p)
		if id != "" {
			venueIDs = append(venueIDs, id)
		}
	}
	if len(venueIDs) == 0 {
		errorResponse(w, http.StatusBadRequest, nil)
		return
	}

	var from, to *time.Time
	if v := strings.TrimSpace(r.URL.Query().Get("from")); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			errorResponse(w, http.StatusBadRequest, nil)
			return
		}
		tu := t.UTC()
		from = &tu
	}
	if v := strings.TrimSpace(r.URL.Query().Get("to")); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			errorResponse(w, http.StatusBadRequest, nil)
			return
		}
		tu := t.UTC()
		to = &tu
	}

	rep, err := h.bookingUsecase.Stats(r.Context(), &bookingModel.StatsReq{
		VenueIds: venueIDs,
		From:     from,
		To:       to,
	})
	if err != nil {
		errorResponse(w, 0, err)
		return
	}

	bookingsByDay := make([]map[string]any, 0, len(rep.BookingsByDay))
	for _, d := range rep.BookingsByDay {
		bookingsByDay = append(bookingsByDay, map[string]any{"date": d.Date, "value": d.Value})
	}
	revenueByDay := make([]map[string]any, 0, len(rep.RevenueByDay))
	for _, d := range rep.RevenueByDay {
		revenueByDay = append(revenueByDay, map[string]any{"date": d.Date, "value": d.Value})
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"total_bookings":     rep.TotalBookings,
		"confirmed_bookings": rep.ConfirmedBookings,
		"cancelled_bookings": rep.CancelledBookings,
		"total_revenue":      rep.TotalRevenue,
		"currency":           rep.Currency,
		"bookings_by_day":    bookingsByDay,
		"revenue_by_day":     revenueByDay,
	})
}
