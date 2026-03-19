package grpc

import (
	"context"

	"github.com/zhamspace/booking/internal/handler/grpc/dto"
	"github.com/zhamspace/booking/internal/usecase/dict"

	"github.com/zhamspace/booking/pkg/proto/booking_v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Dict struct {
	booking_v1.UnsafeDictServer
	dictUsecase *dict.Usecase
}

func NewDict(dictUsecase *dict.Usecase) *Dict {
	return &Dict{
		dictUsecase: dictUsecase,
	}
}

func (h *Dict) Get(ctx context.Context, _ *emptypb.Empty) (*booking_v1.DictRep, error) {
	result := h.dictUsecase.Get(ctx)

	return dto.EncodeDictMain(result), nil
}
