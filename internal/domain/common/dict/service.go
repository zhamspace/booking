package dict

import (
	"context"

	"github.com/zhamspace/booking/internal/domain/common/dict/model"
)

type Service struct {
	repoData RepoDataI
}

func New(repoData RepoDataI) *Service {
	return &Service{
		repoData: repoData,
	}
}

func (s *Service) Get(ctx context.Context) *model.Main {
	return s.repoData.Get(ctx)
}
