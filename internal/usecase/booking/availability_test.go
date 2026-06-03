package booking

import (
	"context"
	"testing"
	"time"

	auditModel "github.com/zhamspace/booking/internal/domain/audit/model"
	bookingModel "github.com/zhamspace/booking/internal/domain/booking/model"
)

type noopAudit struct{}

func (noopAudit) Log(context.Context, *auditModel.Job) {}

type availabilityService struct {
	lastList *bookingModel.ListReq
	total    int64
}

func (s *availabilityService) List(_ context.Context, pars *bookingModel.ListReq) ([]*bookingModel.Main, int64, error) {
	s.lastList = pars
	return nil, s.total, nil
}
func (s *availabilityService) Get(context.Context, *bookingModel.GetReq, bool) (*bookingModel.Main, bool, error) {
	return nil, false, nil
}
func (s *availabilityService) Create(context.Context, *bookingModel.Edit) (string, error) {
	return "", nil
}
func (s *availabilityService) Update(context.Context, *bookingModel.GetReq, *bookingModel.Edit) error {
	return nil
}
func (s *availabilityService) Delete(context.Context, *bookingModel.GetReq) error { return nil }
func (s *availabilityService) Stats(context.Context, *bookingModel.StatsReq) (*bookingModel.StatsRep, error) {
	return nil, nil
}

func TestCountActiveOverlapsUsesActiveAvailabilityFilters(t *testing.T) {
	service := &availabilityService{total: 1}
	uc := New(noopAudit{}, service)
	start := time.Date(2026, 5, 27, 13, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	activeAt := start.Add(-time.Minute)
	excludeSession := "session-1"

	count, err := uc.countActiveOverlaps(context.Background(), "venue-1", "resource-1", &start, &end, activeAt, &excludeSession)
	if err != nil {
		t.Fatalf("countActiveOverlaps returned error: %v", err)
	}
	if count != 1 {
		t.Fatalf("count = %d, want 1", count)
	}
	if service.lastList == nil {
		t.Fatal("service List was not called")
	}
	if service.lastList.ActiveAt == nil || !service.lastList.ActiveAt.Equal(activeAt) {
		t.Fatalf("ActiveAt = %v, want %v", service.lastList.ActiveAt, activeAt)
	}
	if service.lastList.ExcludeSessionID == nil || *service.lastList.ExcludeSessionID != excludeSession {
		t.Fatalf("ExcludeSessionID = %v, want %s", service.lastList.ExcludeSessionID, excludeSession)
	}
	if !service.lastList.OnlyCount {
		t.Fatal("availability overlap check must request only count")
	}
}
