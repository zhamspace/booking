package booking

import (
	"context"

	"github.com/zhamspace/booking/internal/domain/booking/model"
	"github.com/zhamspace/booking/internal/errs"
)

type Service struct {
	repoDb RepoDbI
}

func New(repoDb RepoDbI) *Service {
	return &Service{
		repoDb: repoDb,
	}
}

func (s *Service) List(ctx context.Context, pars *model.ListReq) ([]*model.Main, int64, error) {
	items, totalCount, err := s.repoDb.List(ctx, pars)
	if err != nil {
		return nil, 0, err
	}
	return items, totalCount, nil
}

func (s *Service) Get(ctx context.Context, pars *model.GetReq, errNE bool) (*model.Main, bool, error) {
	item, found, err := s.repoDb.Get(ctx, pars)
	if err != nil {
		return nil, false, err
	}
	if !found {
		if errNE {
			return nil, false, errs.ObjectNotFound
		}
		return nil, false, nil
	}
	return item, found, nil
}

func (s *Service) Create(ctx context.Context, obj *model.Edit) (string, error) {
	return s.repoDb.Create(ctx, obj)
}

func (s *Service) Update(ctx context.Context, pars *model.GetReq, obj *model.Edit) error {
	return s.repoDb.Update(ctx, pars, obj)
}

func (s *Service) Delete(ctx context.Context, pars *model.GetReq) error {
	return s.repoDb.Delete(ctx, pars)
}
