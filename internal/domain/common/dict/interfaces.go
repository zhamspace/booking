package dict

import (
	"context"

	"github.com/zhamspace/booking/internal/domain/common/dict/model"
)

type RepoDataI interface {
	Get(ctx context.Context) *model.Main
}
