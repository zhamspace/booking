package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/zhamspace/booking/internal/domain/booking/model"
	repoModel "github.com/zhamspace/booking/internal/domain/booking/repo/db/model"
	commonRepoPg "github.com/zhamspace/booking/internal/domain/common/repo/pg"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/mechta-market/mobone/v2"
	moboneTools "github.com/mechta-market/mobone/v2/tools"
)

type Repo struct {
	*commonRepoPg.Base
	ModelStore *mobone.ModelStore
	tracer     trace.Tracer
}

func New(con *pgxpool.Pool) *Repo {
	base := commonRepoPg.NewBase(con)
	return &Repo{
		Base: base,
		ModelStore: &mobone.ModelStore{
			Con:       con,
			QB:        base.QB,
			TableName: "booking",
		},
		tracer: otel.Tracer("booking/booking-repo"),
	}
}

func (r *Repo) List(ctx context.Context, pars *model.ListReq) (_ []*model.Main, _ int64, finalError error) {
	ctx, span := r.tracer.Start(ctx, "booking.repo.DB.List")
	defer func() {
		if finalError != nil {
			span.RecordError(finalError)
			span.SetStatus(codes.Error, finalError.Error())
		}
		span.End()
	}()

	conditions, conditionExps := r.getConditions(pars)
	sort := moboneTools.ConstructSortColumns(allowedSortFields, pars.Sort)

	items := make([]*repoModel.Select, 0, len(conditions))
	totalCount, err := r.ModelStore.List(ctx, mobone.ListParams{
		Conditions:           conditions,
		ConditionExpressions: conditionExps,
		Page:                 pars.Page,
		PageSize:             pars.PageSize,
		WithTotalCount:       pars.WithTotalCount,
		OnlyCount:            pars.OnlyCount,
		Sort:                 sort,
	}, func(add bool) mobone.ListModelI {
		item := &repoModel.Select{}
		if add {
			items = append(items, item)
		}
		return item
	})
	if err != nil {
		return nil, 0, fmt.Errorf("ModelStore.List: %w", err)
	}

	return lo.Map(items, repoModel.EncodeMain), totalCount, nil
}

func (r *Repo) Get(ctx context.Context, pars *model.GetReq) (_ *model.Main, _ bool, finalError error) {
	ctx, span := r.tracer.Start(ctx, "booking.repo.DB.Get")
	defer func() {
		if finalError != nil {
			span.RecordError(finalError)
			span.SetStatus(codes.Error, finalError.Error())
		}
		span.End()
	}()

	m := &repoModel.Select{
		Id: pars.Id,
	}

	found, err := r.ModelStore.Get(ctx, m)
	if err != nil {
		return nil, false, fmt.Errorf("ModelStore.Get: %w", err)
	}
	if !found {
		return nil, false, nil
	}

	return repoModel.EncodeMain(m, 0), true, nil
}

func (r *Repo) Create(ctx context.Context, obj *model.Edit) (_ string, finalError error) {
	ctx, span := r.tracer.Start(ctx, "booking.repo.DB.Create")
	defer func() {
		if finalError != nil {
			span.RecordError(finalError)
			span.SetStatus(codes.Error, finalError.Error())
		}
		span.End()
	}()

	m := repoModel.DecodeEdit(obj)

	err := r.ModelStore.Create(ctx, m)
	if err != nil {
		return "", fmt.Errorf("ModelStore.Create: %w", err)
	}

	return m.PKId, nil
}

func (r *Repo) Update(ctx context.Context, pars *model.GetReq, obj *model.Edit) (finalError error) {
	ctx, span := r.tracer.Start(ctx, "booking.repo.DB.Update")
	defer func() {
		if finalError != nil {
			span.RecordError(finalError)
			span.SetStatus(codes.Error, finalError.Error())
		}
		span.End()
	}()

	m := repoModel.DecodeEdit(obj)
	m.PKId = pars.Id
	if m.UpdatedAt == nil {
		now := time.Now().UTC()
		m.UpdatedAt = &now
	}

	err := r.ModelStore.Update(ctx, m)
	if err != nil {
		return fmt.Errorf("ModelStore.Update: %w", err)
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, pars *model.GetReq) (finalError error) {
	ctx, span := r.tracer.Start(ctx, "booking.repo.DB.Delete")
	defer func() {
		if finalError != nil {
			span.RecordError(finalError)
			span.SetStatus(codes.Error, finalError.Error())
		}
		span.End()
	}()

	err := r.ModelStore.Delete(ctx, &repoModel.Upsert{PKId: pars.Id})
	if err != nil {
		return fmt.Errorf("ModelStore.Delete: %w", err)
	}

	return nil
}

func (r *Repo) Stats(ctx context.Context, req *model.StatsReq) (_ *model.StatsRep, finalError error) {
	ctx, span := r.tracer.Start(ctx, "booking.repo.DB.Stats")
	defer func() {
		if finalError != nil {
			span.RecordError(finalError)
			span.SetStatus(codes.Error, finalError.Error())
		}
		span.End()
	}()

	if len(req.VenueIds) == 0 {
		return &model.StatsRep{}, nil
	}

	args := []any{req.VenueIds}
	fromCond := ""
	toCond := ""
	if req.From != nil {
		args = append(args, *req.From)
		fromCond = fmt.Sprintf(" AND start_at >= $%d", len(args))
	}
	if req.To != nil {
		args = append(args, *req.To)
		toCond = fmt.Sprintf(" AND start_at < $%d", len(args))
	}

	query := fmt.Sprintf(`
		SELECT
			COUNT(*) FILTER (WHERE status != 'cancelled') AS total_bookings,
			COUNT(*) FILTER (WHERE status IN ('confirmed','completed')) AS confirmed_bookings,
			COUNT(*) FILTER (WHERE status = 'cancelled') AS cancelled_bookings,
			COALESCE(SUM(price_total) FILTER (WHERE status IN ('confirmed','completed')), 0) AS total_revenue,
			COALESCE(MAX(currency) FILTER (WHERE status IN ('confirmed','completed')), '') AS currency
		FROM booking
		WHERE venue_id = ANY($1) %s%s`, fromCond, toCond)

	rep := &model.StatsRep{}
	err := r.Base.Con.QueryRow(ctx, query, args...).Scan(
		&rep.TotalBookings,
		&rep.ConfirmedBookings,
		&rep.CancelledBookings,
		&rep.TotalRevenue,
		&rep.Currency,
	)
	if err != nil {
		return nil, fmt.Errorf("stats query: %w", err)
	}

	dayArgs := []any{req.VenueIds}
	dayFromCond := ""
	dayToCond := ""
	if req.From != nil {
		dayArgs = append(dayArgs, *req.From)
		dayFromCond = fmt.Sprintf(" AND start_at >= $%d", len(dayArgs))
	}
	if req.To != nil {
		dayArgs = append(dayArgs, *req.To)
		dayToCond = fmt.Sprintf(" AND start_at < $%d", len(dayArgs))
	}

	dayQuery := fmt.Sprintf(`
		SELECT
			start_at::date::text AS date,
			COUNT(*) FILTER (WHERE status != 'cancelled') AS bookings,
			COALESCE(SUM(price_total) FILTER (WHERE status IN ('confirmed','completed')), 0) AS revenue
		FROM booking
		WHERE venue_id = ANY($1) %s%s
		GROUP BY start_at::date
		ORDER BY start_at::date`, dayFromCond, dayToCond)

	rows, err := r.Base.Con.Query(ctx, dayQuery, dayArgs...)
	if err != nil {
		return nil, fmt.Errorf("day stats query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var date string
		var bookingCount, revenue int64
		if err := rows.Scan(&date, &bookingCount, &revenue); err != nil {
			continue
		}
		rep.BookingsByDay = append(rep.BookingsByDay, model.DayStat{Date: date, Value: bookingCount})
		rep.RevenueByDay = append(rep.RevenueByDay, model.DayStat{Date: date, Value: revenue})
	}

	return rep, nil
}
