package booking

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/zhamspace/booking/internal/constant"
	auditConstants "github.com/zhamspace/booking/internal/domain/audit/constants"
	auditModel "github.com/zhamspace/booking/internal/domain/audit/model"
	bookingConstants "github.com/zhamspace/booking/internal/domain/booking/constants"
	bookingModel "github.com/zhamspace/booking/internal/domain/booking/model"
	commonModel "github.com/zhamspace/booking/internal/domain/common/model"
	"github.com/zhamspace/booking/internal/domain/common/session"
	"github.com/zhamspace/booking/internal/domain/common/tools"
	"github.com/zhamspace/booking/internal/errs"
	"google.golang.org/grpc/codes"
)

type Usecase struct {
	auditService   AuditServiceI
	bookingService ServiceI
}

func New(
	auditService AuditServiceI,
	bookingService ServiceI,
) *Usecase {
	return &Usecase{
		auditService:   auditService,
		bookingService: bookingService,
	}
}

func (u Usecase) List(ctx context.Context, pars *bookingModel.ListReq) (_ []*bookingModel.Main, _ int64, finalError error) {
	if pars == nil {
		return nil, 0, errs.New(errs.InvalidConfig, "list params are required", nil)
	}

	fields := map[string]string{}
	if !pars.OnlyCount {
		if err := tools.RequirePageSize(pars.ListParams, constant.MaxPageSize); err != nil {
			fields["list_params.page_size"] = fmt.Sprintf("must be between 1 and %d", constant.MaxPageSize)
		}
	}
	if len(fields) > 0 {
		return nil, 0, errs.Validation(fields)
	}
	if pars.From != nil && pars.To != nil && pars.From.After(*pars.To) {
		return nil, 0, errs.Validation(map[string]string{
			"to": "must be greater than or equal to from",
		})
	}

	items, totalCount, err := u.bookingService.List(ctx, pars)
	if err != nil {
		fields["method"] = "bookingService.List"
		return nil, 0, errs.Wrap(errs.ServiceNA, "failed data retrieval", fields, err)
	}

	return items, totalCount, nil
}

func (u Usecase) Get(ctx context.Context, pars *bookingModel.GetReq) (*bookingModel.Main, error) {
	fields := map[string]string{}

	item, _, err := u.bookingService.Get(ctx, pars, true)
	if err != nil {
		fields["method"] = "bookingService.Get"
		return nil, errs.Wrap(errs.ServiceNA, "failed data retrieval", fields, err)
	}

	return item, nil
}

func (u Usecase) Create(ctx context.Context, obj *bookingModel.Edit) (*bookingModel.Main, error) {
	if obj == nil {
		return nil, errs.New(errs.InvalidConfig, "booking payload is required", nil)
	}

	fields := map[string]string{}
	now := time.Now().UTC()

	venueID := normalizeRequiredString(obj.VenueId)
	if venueID == "" {
		fields["venue_id"] = "required"
	}

	resourceID := normalizeRequiredString(obj.ResourceId)
	if resourceID == "" {
		fields["resource_id"] = "required"
	}

	timezone := normalizeRequiredString(obj.Timezone)
	if timezone == "" {
		fields["timezone"] = "required"
	}

	var userID *string
	if obj.UserId != nil {
		userID = normalizeOptionalString(obj.UserId)
	}
	if userID == nil {
		sessUserID := strings.TrimSpace(session.ExtractFromContext(ctx).Sub)
		if sessUserID != "" {
			userID = &sessUserID
		}
	}
	if userID == nil {
		fields["user_id"] = "required"
	}

	if obj.StartAt == nil {
		fields["start_at"] = "required"
	}
	if obj.EndAt == nil {
		fields["end_at"] = "required"
	}
	if obj.StartAt != nil && obj.EndAt != nil && !obj.EndAt.After(*obj.StartAt) {
		fields["end_at"] = "must be greater than start_at"
	}

	priceTotal := obj.PriceTotal
	if priceTotal == nil {
		defaultPrice := int64(0)
		priceTotal = &defaultPrice
	}
	if priceTotal != nil && *priceTotal < 0 {
		fields["price_total"] = "must be greater than or equal to 0"
	}
	if obj.StartAt != nil && obj.StartAt.Before(now) {
		fields["start_at"] = "must be in the future"
	}

	currency := normalizeOptionalString(obj.Currency)
	if currency == nil {
		defaultCurrency := bookingConstants.DefaultCurrency
		currency = &defaultCurrency
	}

	if len(fields) > 0 {
		return nil, errs.Validation(fields)
	}

	overlapCount, err := u.countActiveOverlaps(ctx, venueID, resourceID, obj.StartAt, obj.EndAt, now)
	if err != nil {
		fields["method"] = "bookingService.List"
		return nil, errs.Wrap(errs.ServiceNA, "failed availability check", fields, err)
	}
	if overlapCount >= bookingConstants.DefaultCapacity {
		return nil, errs.WithStatus(
			errs.InvalidRequest,
			"booking time slot is unavailable",
			map[string]string{
				"venue_id":    venueID,
				"resource_id": resourceID,
				"start_at":    obj.StartAt.Format(time.RFC3339),
				"end_at":      obj.EndAt.Format(time.RFC3339),
				"reason":      "resource capacity exceeded for selected time range",
			},
			codes.Aborted,
			http.StatusConflict,
		)
	}

	statusCreated := bookingConstants.StatusCreated
	paymentStatusPending := bookingConstants.PaymentStatusPending
	holdExpiresAt := now.Add(bookingConstants.DefaultHoldDuration)
	id, err := u.bookingService.Create(ctx, &bookingModel.Edit{
		VenueId:         &venueID,
		ResourceId:      &resourceID,
		SessionId:       normalizeOptionalString(obj.SessionId),
		UserId:          userID,
		CreatedAt:       &now,
		UpdatedAt:       &now,
		Status:          &statusCreated,
		PaymentStatus:   &paymentStatusPending,
		PriceTotal:      priceTotal,
		Currency:        currency,
		StartAt:         obj.StartAt,
		EndAt:           obj.EndAt,
		Timezone:        &timezone,
		HoldExpiresAt:   &holdExpiresAt,
		CancelReason:    normalizeOptionalString(obj.CancelReason),
		PaymentIntentId: normalizeOptionalString(obj.PaymentIntentId),
	})
	if err != nil {
		fields["method"] = "bookingService.Create"
		return nil, errs.Wrap(errs.ServiceNA, "failed create booking", fields, err)
	}

	item, _, err := u.bookingService.Get(ctx, &bookingModel.GetReq{Id: id}, true)
	if err != nil {
		fields["method"] = "bookingService.Get"
		return nil, errs.Wrap(errs.ServiceNA, "failed data retrieval", fields, err)
	}

	u.auditService.Log(ctx, &auditModel.Job{
		ChangeType: auditConstants.ChangeTypeCreated,
		Resource:   auditConstants.ResourceBooking,
		ObjectId:   item.Id,
		Object:     item,
	})

	return item, nil
}

func (u Usecase) countActiveOverlaps(ctx context.Context, venueID, resourceID string, startAt, endAt *time.Time, activeAt time.Time) (int64, error) {
	_, totalCount, err := u.bookingService.List(ctx, &bookingModel.ListReq{
		ListParams: commonModel.ListParams{
			OnlyCount: true,
		},
		VenueId:    &venueID,
		ResourceId: &resourceID,
		From:       startAt,
		To:         endAt,
		ActiveAt:   &activeAt,
	})
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

func normalizeRequiredString(v *string) string {
	if v == nil {
		return ""
	}

	return strings.TrimSpace(*v)
}

func normalizeOptionalString(v *string) *string {
	if v == nil {
		return nil
	}

	s := strings.TrimSpace(*v)
	if s == "" {
		return nil
	}

	return &s
}
