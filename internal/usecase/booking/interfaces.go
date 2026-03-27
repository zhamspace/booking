package booking

import (
	"context"

	auditModel "github.com/zhamspace/booking/internal/domain/audit/model"
	bookingModel "github.com/zhamspace/booking/internal/domain/booking/model"
)

type AuditServiceI interface {
	Log(ctx context.Context, job *auditModel.Job)
}

type ServiceI interface {
	List(ctx context.Context, pars *bookingModel.ListReq) ([]*bookingModel.Main, int64, error)
	Get(ctx context.Context, pars *bookingModel.GetReq, errNE bool) (*bookingModel.Main, bool, error)
	Create(ctx context.Context, obj *bookingModel.Edit) (string, error)
	Update(ctx context.Context, pars *bookingModel.GetReq, obj *bookingModel.Edit) error
	Delete(ctx context.Context, pars *bookingModel.GetReq) error
}
