package dict

import (
	"context"

	dictModel "github.com/zhamspace/booking/internal/domain/common/dict/model"
)

type Usecase struct {
	dictService ServiceI
}

func New(dictService ServiceI) *Usecase {
	return &Usecase{
		dictService: dictService,
	}
}

func (u *Usecase) Get(ctx context.Context) *dictModel.Main {
	return u.dictService.Get(ctx)
}
