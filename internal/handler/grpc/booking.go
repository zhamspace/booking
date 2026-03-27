package grpc

import (
	"context"

	"github.com/zhamspace/booking/internal/errs"
	"github.com/zhamspace/booking/internal/handler/grpc/dto"

	bookingModel "github.com/zhamspace/booking/internal/domain/booking/model"
	usecaseBookingP "github.com/zhamspace/booking/internal/usecase/booking"

	"github.com/zhamspace/booking/pkg/proto/booking_v1"
	"github.com/zhamspace/booking/pkg/proto/common"

	"github.com/samber/lo"
)

type Booking struct {
	booking_v1.UnsafeBookingServer
	bookingUsecase *usecaseBookingP.Usecase
}

func NewBooking(bookingUsecase *usecaseBookingP.Usecase) *Booking {
	return &Booking{
		bookingUsecase: bookingUsecase,
	}
}

func (h *Booking) List(ctx context.Context, req *booking_v1.BookingListReq) (*booking_v1.BookingListRep, error) {
	if req.ListParams == nil {
		req.ListParams = &common.ListParamsSt{}
	}

	items, totalCount, err := h.bookingUsecase.List(ctx, &bookingModel.ListReq{
		ListParams:       dto.DecodeListParams(req.ListParams),
		Ids:              dto.ToSlicePtr(req.Ids),
		VenueId:          dto.DecodeOptionalString(req.VenueId),
		ResourceId:       dto.DecodeOptionalString(req.ResourceId),
		UserId:           dto.DecodeOptionalString(req.UserId),
		SessionId:        dto.DecodeOptionalString(req.SessionId),
		Status:           dto.DecodeOptionalString(req.Status),
		PaymentStatus:    dto.DecodeOptionalString(req.PaymentStatus),
		From:             dto.DecodeTimestampPtr(req.From),
		To:               dto.DecodeTimestampPtr(req.To),
		IncludeCancelled: req.IncludeCancelled,
	})
	if err != nil {
		return nil, err
	}

	return &booking_v1.BookingListRep{
		PaginationInfo: dto.EncodePaginationInfoSt(req.ListParams, totalCount),
		Results:        lo.Map(items, dto.EncodeBookingMain),
	}, nil
}

func (h *Booking) Get(ctx context.Context, req *booking_v1.BookingGetReq) (*booking_v1.BookingMain, error) {
	item, err := h.bookingUsecase.Get(ctx, &bookingModel.GetReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return dto.EncodeBookingMain(item, 0), nil
}

func (h *Booking) Create(ctx context.Context, req *booking_v1.BookingCreateReq) (*booking_v1.BookingMain, error) {
	item, err := h.bookingUsecase.Create(ctx, &bookingModel.Edit{
		VenueId:    new(req.VenueId),
		ResourceId: new(req.ResourceId),
		SessionId:  dto.DecodeOptionalString(req.SessionId),
		UserId:     dto.DecodeOptionalString(req.UserId),
		PriceTotal: req.PriceTotal,
		Currency:   dto.DecodeOptionalString(req.Currency),
		StartAt:    dto.DecodeTimestampPtr(req.StartAt),
		EndAt:      dto.DecodeTimestampPtr(req.EndAt),
		Timezone:   new(req.Timezone),
	})
	if err != nil {
		return nil, err
	}

	return dto.EncodeBookingMain(item, 0), nil
}

func (h *Booking) Confirm(ctx context.Context, req *booking_v1.BookingConfirmReq) (*booking_v1.BookingMain, error) {
	return nil, errs.NotImplemented
}

func (h *Booking) Cancel(ctx context.Context, req *booking_v1.BookingCancelReq) (*booking_v1.BookingMain, error) {
	return nil, errs.NotImplemented
}
