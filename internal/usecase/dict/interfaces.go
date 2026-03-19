package dict

import (
	"context"

	dictModel "github.com/zhamspace/booking/internal/domain/common/dict/model"
)

type ServiceI interface {
	Get(ctx context.Context) *dictModel.Main
}
