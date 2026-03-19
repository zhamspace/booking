package grpc

import (
	"context"

	"github.com/zhamspace/booking/internal/handler/grpc/dto"

	auditDomain "github.com/zhamspace/booking/internal/domain/audit"
	auditModel "github.com/zhamspace/booking/internal/domain/audit/model"
	usecaseAuditP "github.com/zhamspace/booking/internal/usecase/audit"

	"github.com/zhamspace/booking/pkg/proto/booking_v1"
	"github.com/zhamspace/booking/pkg/proto/common"

	"github.com/samber/lo"
)

type Audit struct {
	booking_v1.UnsafeAuditServer
	auditUsecase *usecaseAuditP.Usecase
	serializer   auditDomain.ObjectSerializer
}

func NewAudit(auditUsecase *usecaseAuditP.Usecase, serializer auditDomain.ObjectSerializer) *Audit {
	return &Audit{
		auditUsecase: auditUsecase,
		serializer:   serializer,
	}
}

func (h *Audit) List(ctx context.Context, req *booking_v1.AuditListReq) (*booking_v1.AuditListRep, error) {
	if req.ListParams == nil {
		req.ListParams = &common.ListParamsSt{}
	}

	items, totalCount, err := h.auditUsecase.List(ctx, &auditModel.ListReq{
		ListParams: dto.DecodeListParams(req.ListParams),
		Ids:        dto.ToSlicePtr(req.Ids),
		UserId:     req.UserId,
		ObjectId:   req.ObjectId,
		Resource:   req.Resource,
		ChangeType: req.ChangeType,
	})
	if err != nil {
		return nil, err
	}

	encoder := dto.AuditEncoder{Serializer: h.serializer}
	return &booking_v1.AuditListRep{
		PaginationInfo: dto.EncodePaginationInfoSt(req.ListParams, totalCount),
		Results:        lo.Map(items, encoder.EncodeAuditMain),
	}, nil
}

func (h *Audit) Get(ctx context.Context, req *booking_v1.AuditGetReq) (*booking_v1.AuditMain, error) {
	item, err := h.auditUsecase.Get(ctx, &auditModel.GetReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	encoder := dto.AuditEncoder{Serializer: h.serializer}
	return encoder.EncodeAuditMain(item, 0), nil
}
