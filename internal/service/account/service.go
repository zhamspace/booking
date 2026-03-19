package account

import (
	"context"
	"fmt"

	"github.com/zhamspace/booking/internal/errs"
	"github.com/zhamspace/booking/internal/service/account/model"
)

type Service struct {
	repo RepoI
}

func New(repo RepoI) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetUser(ctx context.Context, pars *model.GetReq, errNE bool) (*model.Main, bool, error) {
	if pars == nil {
		pars = &model.GetReq{}
	}
	if *pars == (model.GetReq{}) {
		return nil, false, fmt.Errorf("GetUserById empty request: %w", errs.InvalidRequest)
	}
	item, found, err := s.repo.GetUser(ctx, pars)
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
