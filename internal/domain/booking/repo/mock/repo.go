package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	bookingModel "github.com/zhamspace/booking/internal/domain/booking/model"
)

type Repo struct {
}

func New() *Repo {
	return &Repo{}
}

func (r Repo) List(ctx context.Context, pars *bookingModel.ListReq) ([]*bookingModel.Main, int64, error) {
	const (
		venueId    = "ee70f720-1067-418e-818c-f17d8581e44d"
		resourceId = "5ea59285-3301-4dc7-9684-35dcd477523d"
	)

	gofakeit.Seed(0)

	totalCount := int64(1000) // simulate big dataset

	start := (pars.Page - 1) * pars.PageSize
	end := start + pars.PageSize
	if end > int64(totalCount) {
		end = int64(totalCount)
	}

	items := make([]*bookingModel.Main, 0, pars.PageSize)

	for i := start; i < end; i++ {
		items = append(items, &bookingModel.Main{
			Id: gofakeit.UUID(),
			//VenueId:       gofakeit.UUID(),
			VenueId:    venueId,
			ResourceId: resourceId,
			SessionId:  new(gofakeit.UUID()),
			UserId: gofakeit.RandString(
				[]string{
					"8f1fce7a-2a11-4e5b-a82e-7697e0460fc4",
					"5d62a0c4-2066-4e00-8f77-554eff5c4ff6",
					"26e4cbcf-3b99-4176-9b19-fe2592db175d",
					"3aa0324c-0f5a-45c3-a00b-76e5a4bfe57e",
					"9da9062c-e2b9-40c3-9d06-be6a68282a78",
					"79e8005e-464f-4515-8695-4df5665280b4",
					"ce13f804-79ff-442c-a999-15eba8bfe854",
					"aa8a9c6b-4567-427b-abc3-627a596473db",
					"388428ad-ae63-45c8-8fb5-b1a7c5f9dfb4",
					"c9a19d96-f8ac-4982-be25-9bd26ccf4f9d",
					"1b89fdcc-5d6b-4f2e-be53-e06f1fdeef30",
					"5803c267-77c3-4444-95b6-da743f56afe0",
				}),
			CreatedAt:     gofakeit.DateRange(time.Now().AddDate(0, -1, 0), time.Now()),
			UpdatedAt:     time.Now(),
			Status:        gofakeit.RandString([]string{"created", "confirmed", "cancelled"}),
			PaymentStatus: gofakeit.RandString([]string{"pending", "paid", "failed"}),
			PriceTotal:    int64(gofakeit.Price(2, 40)) * 250,
			Currency:      "KZT",
			//StartAt:         new(gofakeit.DateRange(time.Now().AddDate(0, 0, 0), time.Now()).AddDate(0, 0, 3)),
			//EndAt:           new(gofakeit.DateRange(time.Now().AddDate(0, 0, 0), time.Now()).AddDate(0, 0, 3)),
			Timezone:        "Asia/Almaty",
			HoldExpiresAt:   new(gofakeit.DateRange(time.Now(), time.Now().AddDate(0, 0, 7))),
			PaymentIntentId: new(gofakeit.UUID()),
			CancelReason:    nil,
			CancelledAt:     nil,
			ConfirmedAt:     new(gofakeit.DateRange(time.Now(), time.Now().AddDate(0, 0, 7))),
			CompletedAt:     nil,
		})
	}

	return items, totalCount, nil
}

func (r Repo) Get(ctx context.Context, pars *bookingModel.GetReq) (*bookingModel.Main, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repo) Create(ctx context.Context, obj *bookingModel.Edit) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repo) Update(ctx context.Context, pars *bookingModel.GetReq, obj *bookingModel.Edit) error {
	//TODO implement me
	panic("implement me")
}

func (r Repo) Delete(ctx context.Context, pars *bookingModel.GetReq) error {
	//TODO implement me
	panic("implement me")
}

func (r Repo) Stats(ctx context.Context, req *bookingModel.StatsReq) (*bookingModel.StatsRep, error) {
	//TODO implement me
	panic("implement me")
}
