package http

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	bookingModel "github.com/zhamspace/booking/internal/domain/booking/model"
	commonModel "github.com/zhamspace/booking/internal/domain/common/model"
	commonSession "github.com/zhamspace/booking/internal/domain/common/session"
	sessionClient "github.com/zhamspace/booking/internal/service/session"
	usecaseBookingP "github.com/zhamspace/booking/internal/usecase/booking"
)

type ReservationType string

const (
	reservationTypeAll      ReservationType = "all"
	reservationTypeBookings ReservationType = "bookings"
	reservationTypeSessions ReservationType = "sessions"
)

type Reservation struct {
	ItemType          string     `json:"item_type"`
	ID                string     `json:"id"`
	VenueID           string     `json:"venue_id"`
	ResourceID        string     `json:"resource_id,omitempty"`
	SessionID         *string    `json:"session_id,omitempty"`
	UserID            string     `json:"user_id,omitempty"`
	Name              string     `json:"name,omitempty"`
	Status            string     `json:"status"`
	ParticipantStatus string     `json:"participant_status,omitempty"`
	PaymentStatus     string     `json:"payment_status,omitempty"`
	PriceTotal        int64      `json:"price_total"`
	Currency          string     `json:"currency,omitempty"`
	StartAt           *time.Time `json:"start_at,omitempty"`
	EndAt             *time.Time `json:"end_at,omitempty"`
	Timezone          string     `json:"timezone,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	HoldExpiresAt     *time.Time `json:"hold_expires_at,omitempty"`
	PaymentIntentID   *string    `json:"payment_intent_id,omitempty"`
	CancelReason      *string    `json:"cancel_reason,omitempty"`
	CancelledAt       *time.Time `json:"cancelled_at,omitempty"`
	ConfirmedAt       *time.Time `json:"confirmed_at,omitempty"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
}

type Reservations struct {
	bookingUsecase *usecaseBookingP.Usecase
	sessionClient  *sessionClient.Client
}

func NewReservations(bookingUsecase *usecaseBookingP.Usecase, sessionClient *sessionClient.Client) *Reservations {
	return &Reservations{bookingUsecase: bookingUsecase, sessionClient: sessionClient}
}

func (h *Reservations) List(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	q := r.URL.Query()
	page := parseInt64(q.Get("list_params.page"), 0)
	pageSize := parseInt64(q.Get("list_params.page_size"), 30)
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 30
	}
	sortParam := strings.TrimSpace(q.Get("list_params.sort"))
	if sortParam == "" {
		sortParam = "-end_at"
	}

	userID := strings.TrimSpace(q.Get("user_id"))
	if userID == "" {
		userID = strings.TrimSpace(commonSession.ExtractFromContext(r.Context()).Sub)
	}
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "authentication required"})
		return
	}

	kind := ReservationType(strings.TrimSpace(strings.ToLower(q.Get("type"))))
	if kind == "" {
		kind = reservationTypeAll
	}
	if kind != reservationTypeAll && kind != reservationTypeBookings && kind != reservationTypeSessions {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid type"})
		return
	}

	includeCancelled := true
	if raw := strings.TrimSpace(q.Get("include_cancelled")); raw != "" {
		if v, err := strconv.ParseBool(raw); err == nil {
			includeCancelled = v
		}
	}
	var from, to *time.Time
	if parsed, ok := parseTimeQuery(q.Get("from")); ok {
		from = parsed
	} else if strings.TrimSpace(q.Get("from")) != "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid from"})
		return
	}
	if parsed, ok := parseTimeQuery(q.Get("to")); ok {
		to = parsed
	} else if strings.TrimSpace(q.Get("to")) != "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid to"})
		return
	}

	fetchPage := page
	fetchSize := pageSize
	if kind == reservationTypeAll {
		fetchPage = 0
		fetchSize = (page + 1) * pageSize
	}

	var rows []Reservation
	var total int64
	if kind == reservationTypeAll || kind == reservationTypeBookings {
		items, count, err := h.bookingUsecase.List(r.Context(), &bookingModel.ListReq{
			ListParams:       bookingListParams(fetchPage, fetchSize, sortParam),
			UserId:           &userID,
			From:             from,
			To:               to,
			IncludeCancelled: includeCancelled,
		})
		if err != nil {
			errorResponse(w, 0, err)
			return
		}
		total += count
		for _, item := range items {
			rows = append(rows, reservationFromBooking(item))
		}
	}

	if kind == reservationTypeAll || kind == reservationTypeSessions {
		rep, err := h.sessionClient.ListParticipated(r.Context(), sessionClient.ListParticipatedReq{
			UserID:           userID,
			Page:             fetchPage,
			PageSize:         fetchSize,
			Sort:             sortParam,
			From:             from,
			To:               to,
			IncludeCancelled: includeCancelled,
		})
		if err != nil {
			errorResponse(w, http.StatusBadGateway, err)
			return
		}
		total += rep.PaginationInfo.TotalCount
		for _, item := range rep.Results {
			rows = append(rows, reservationFromSession(item))
		}
	}

	sortReservations(rows, sortParam)
	if kind == reservationTypeAll {
		start := int(page * pageSize)
		if start > len(rows) {
			rows = nil
		} else {
			end := start + int(pageSize)
			if end > len(rows) {
				end = len(rows)
			}
			rows = rows[start:end]
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"results": rows,
		"pagination_info": map[string]any{
			"page":        page,
			"page_size":   pageSize,
			"total_count": total,
		},
	})
}

func reservationFromBooking(b *bookingModel.Main) Reservation {
	return Reservation{
		ItemType:        "booking",
		ID:              b.Id,
		VenueID:         b.VenueId,
		ResourceID:      b.ResourceId,
		SessionID:       b.SessionId,
		UserID:          b.UserId,
		Status:          b.Status,
		PaymentStatus:   b.PaymentStatus,
		PriceTotal:      b.PriceTotal,
		Currency:        b.Currency,
		StartAt:         b.StartAt,
		EndAt:           b.EndAt,
		Timezone:        b.Timezone,
		CreatedAt:       &b.CreatedAt,
		UpdatedAt:       &b.UpdatedAt,
		HoldExpiresAt:   b.HoldExpiresAt,
		PaymentIntentID: b.PaymentIntentId,
		CancelReason:    b.CancelReason,
		CancelledAt:     b.CancelledAt,
		ConfirmedAt:     b.ConfirmedAt,
		CompletedAt:     b.CompletedAt,
	}
}

func reservationFromSession(s sessionClient.ParticipatedSession) Reservation {
	sessionID := s.SessionID
	if sessionID == "" {
		sessionID = s.ID
	}
	return Reservation{
		ItemType:          "session",
		ID:                s.ID,
		VenueID:           s.VenueID,
		ResourceID:        s.ResourceID,
		SessionID:         &sessionID,
		UserID:            s.UserID,
		Name:              s.Name,
		Status:            s.Status,
		ParticipantStatus: s.ParticipantStatus,
		PaymentStatus:     s.PaymentStatus,
		PriceTotal:        s.PriceTotal,
		Currency:          s.Currency,
		StartAt:           s.StartAt,
		EndAt:             s.EndAt,
		Timezone:          s.Timezone,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
		CancelledAt:       s.CancelledAt,
		PaymentIntentID:   s.PaymentIntentID,
	}
}

func bookingListParams(page, pageSize int64, sortParam string) commonModel.ListParams {
	return commonModel.ListParams{
		Page:           page,
		PageSize:       pageSize,
		WithTotalCount: true,
		Sort:           []string{sortParam},
	}
}

func sortReservations(rows []Reservation, sortParam string) {
	desc := strings.HasPrefix(sortParam, "-")
	field := strings.TrimPrefix(sortParam, "-")
	if field == "" {
		field = "end_at"
	}
	sort.SliceStable(rows, func(i, j int) bool {
		left := reservationSortTime(rows[i], field)
		right := reservationSortTime(rows[j], field)
		if desc {
			return left.After(right)
		}
		return left.Before(right)
	})
}

func reservationSortTime(r Reservation, field string) time.Time {
	switch field {
	case "start_at", "starts_at":
		if r.StartAt != nil {
			return *r.StartAt
		}
	case "created_at":
		if r.CreatedAt != nil {
			return *r.CreatedAt
		}
	default:
		if r.EndAt != nil {
			return *r.EndAt
		}
	}
	return time.Time{}
}

func parseTimeQuery(raw string) (*time.Time, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, true
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, false
	}
	tu := t.UTC()
	return &tu, true
}

func parseInt64(raw string, fallback int64) int64 {
	if raw == "" {
		return fallback
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return fallback
	}
	return v
}
