package booking

import (
	"context"

	"github.com/zhamspace/booking/internal/domain/booking/model"
)

type RepoDbI interface {
	List(ctx context.Context, pars *model.ListReq) ([]*model.Main, int64, error)
	Get(ctx context.Context, pars *model.GetReq) (*model.Main, bool, error)
	Create(ctx context.Context, obj *model.Edit) (string, error)
	Update(ctx context.Context, pars *model.GetReq, obj *model.Edit) error
	Delete(ctx context.Context, pars *model.GetReq) error
}
