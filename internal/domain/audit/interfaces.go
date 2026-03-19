package audit

import (
	"context"
	"encoding/json"

	"github.com/zhamspace/booking/internal/domain/audit/model"
)

type RepoDbI interface {
	List(ctx context.Context, pars *model.ListReq) ([]*model.Main, int64, error)
	Get(ctx context.Context, pars *model.GetReq) (*model.Main, bool, error)
	Create(ctx context.Context, obj *model.Edit) (string, error)
}

type ObjectSerializer interface {
	Serialize(resource string, object any) (json.RawMessage, error)
	Deserialize(resource string, data json.RawMessage) (any, error)
}
