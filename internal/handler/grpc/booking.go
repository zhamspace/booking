package grpc

import (
	"context"

	"github.com/zhamspace/booking/internal/errs"
	//"github.com/zhamspace/booking/internal/handler/grpc/dto"

	usecaseBookingP "github.com/zhamspace/booking/internal/usecase/booking"

	"github.com/zhamspace/booking/pkg/proto/booking_v1"
	"github.com/zhamspace/booking/pkg/proto/common"
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
	return nil, errs.NotImplemented

}

func (h *Booking) Get(ctx context.Context, req *booking_v1.BookingGetReq) (*booking_v1.BookingMain, error) {
	return nil, errs.NotImplemented

}

func (h *Booking) Create(ctx context.Context, req *booking_v1.BookingCreateReq) (*booking_v1.BookingMain, error) {
	return nil, errs.NotImplemented
}

func (h *Booking) Confirm(ctx context.Context, req *booking_v1.BookingConfirmReq) (*booking_v1.BookingMain, error) {
	return nil, errs.NotImplemented
}

func (h *Booking) Cancel(ctx context.Context, req *booking_v1.BookingCancelReq) (*booking_v1.BookingMain, error) {
	return nil, errs.NotImplemented
}
