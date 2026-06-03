package repo

import (
	"strings"

	bookingConstants "github.com/zhamspace/booking/internal/domain/booking/constants"
	"github.com/zhamspace/booking/internal/domain/booking/model"
)

var allowedSortFields = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
	"status":     "status",
	"start_at":   "start_at",
	"end_at":     "end_at",
}

func (r *Repo) getConditions(pars *model.ListReq) (map[string]any, map[string][]any) {
	conditions := make(map[string]any)
	conditionExps := make(map[string][]any)

	if pars.Ids != nil {
		conditions["id"] = pars.Ids
	}
	if pars.VenueIds != nil && len(*pars.VenueIds) > 0 {
		conditions["venue_id"] = *pars.VenueIds
	} else if pars.VenueId != nil {
		conditions["venue_id"] = *pars.VenueId
	}
	if pars.ResourceId != nil {
		conditions["resource_id"] = *pars.ResourceId
	}
	if pars.UserId != nil {
		conditions["user_id"] = *pars.UserId
	}
	if pars.SessionId != nil {
		conditions["session_id"] = *pars.SessionId
	}
	if pars.ExcludeSessionID != nil && strings.TrimSpace(*pars.ExcludeSessionID) != "" {
		conditionExps["(session_id is null or session_id <> ?)"] = []any{strings.TrimSpace(*pars.ExcludeSessionID)}
	}
	if pars.Status != nil {
		conditions["status"] = *pars.Status
	}
	if pars.PaymentStatus != nil {
		conditions["payment_status"] = *pars.PaymentStatus
	}
	if pars.Currency != nil {
		conditions["currency"] = *pars.Currency
	}
	if pars.From != nil {
		conditionExps["end_at >= ?"] = []any{*pars.From}
	}
	if pars.To != nil {
		conditionExps["start_at <= ?"] = []any{*pars.To}
	}
	if !pars.IncludeCancelled {
		conditionExps["status <> ?"] = []any{bookingConstants.StatusCancelled}
	}
	if pars.ActiveAt != nil {
		conditionExps["((status = ? and (hold_expires_at is null or hold_expires_at > ?)) or status in (?, ?))"] = []any{
			bookingConstants.StatusCreated,
			*pars.ActiveAt,
			bookingConstants.StatusConfirmed,
			bookingConstants.StatusCompleted,
		}
	}

	return conditions, conditionExps
}
