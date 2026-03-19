package app

import (
	"fmt"

	"github.com/zhamspace/booking/internal/config"

	serviceAcocuntP "github.com/zhamspace/booking/internal/service/account"
	serviceAccountRepoGrpcP "github.com/zhamspace/booking/internal/service/account/repo/grpc"
	serviceAccountRepoRestP "github.com/zhamspace/booking/internal/service/account/repo/rest"
)

func NewAccountRepo() (serviceAcocuntP.RepoI, error) {
	cfg := config.Conf

	if cfg.AccountGrpcUrl != "" {
		conn, err := newGrpcClientConn(
			cfg.AccountGrpcUrl,
			cfg.AccountGrpcSecure,
			cfg.AccountGrpcUsername,
			cfg.AccountGrpcPassword,
			"Account: ",
		)
		if err != nil {
			return nil, fmt.Errorf("grpc conn: %w", err)
		}

		repo, err := serviceAccountRepoGrpcP.New(&serviceAccountRepoGrpcP.ConfigSt{
			GrpcClient: conn,
		})
		if err != nil {
			return nil, fmt.Errorf("grpc repo: %w", err)
		}

		return repo, nil
	}

	if cfg.AccountHttpUrl != "" {
		repo, err := serviceAccountRepoRestP.New(&serviceAccountRepoRestP.ConfigSt{
			Uri: cfg.AccountHttpUrl,
		})
		if err != nil {
			return nil, fmt.Errorf("http repo: %w", err)
		}

		return repo, nil
	}

	return nil, fmt.Errorf("no account transport configured")
}
