package audit

import (
	"context"

	auditModel "github.com/zhamspace/booking/internal/domain/audit/model"
	accountModel "github.com/zhamspace/booking/internal/service/account/model"
)

type ServiceI interface {
	List(ctx context.Context, pars *auditModel.ListReq) ([]*auditModel.Main, int64, error)
	Get(ctx context.Context, pars *auditModel.GetReq, errNE bool) (*auditModel.Main, bool, error)
}

type AccountServiceI interface {
	GetUser(ctx context.Context, pars *accountModel.GetReq, errNE bool) (*accountModel.Main, bool, error)
}
