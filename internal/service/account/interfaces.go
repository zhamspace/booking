package account

import (
	"context"

	"github.com/zhamspace/booking/internal/service/account/model"
)

type RepoI interface {
	GetUser(ctx context.Context, pars *model.GetReq) (*model.Main, bool, error)
}
