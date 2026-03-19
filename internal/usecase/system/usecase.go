package system

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/zhamspace/booking/internal/constant"
	"github.com/zhamspace/booking/internal/domain/common/logging"
	"github.com/zhamspace/booking/internal/errs"
	"github.com/zhamspace/booking/internal/usecase/common"

	commandServiceP "github.com/zhamspace/booking/internal/service/system/command"
	migrationServiceP "github.com/zhamspace/booking/internal/service/system/migration"
)

type Usecase struct {
	commandService   *commandServiceP.Service
	migrationService *migrationServiceP.Service
	loggingService   LoggingServiceI
}

func New(
	commandService *commandServiceP.Service,
	migrationService *migrationServiceP.Service,
	loggingService LoggingServiceI,
) *Usecase {
	return &Usecase{
		commandService:   commandService,
		migrationService: migrationService,
		loggingService:   loggingService,
	}
}

func (u *Usecase) Command(ctx context.Context, num int) error {
	if err := common.CheckPermission(ctx, []string{constant.RoleAdmin}); err != nil {
		return err
	}

	err := u.commandService.Run(num)
	if err != nil {
		return errs.Wrap(errs.ServiceNA, "command failed", map[string]string{
			"command": fmt.Sprint(num),
		}, err)
	}

	return nil
}

func (u *Usecase) MigrationUp(ctx context.Context) error {
	if err := common.CheckPermission(ctx, []string{constant.RoleAdmin}); err != nil {
		return err
	}

	err := u.migrationService.Up()
	if err != nil {
		return errs.Wrap(errs.ServiceNA, "migration up failed", nil, err)
	}

	slog.Info("migrationService.Up success")
	return nil
}

func (u *Usecase) MigrationDownOne(ctx context.Context) error {
	if err := common.CheckPermission(ctx, []string{constant.RoleAdmin}); err != nil {
		return err
	}

	err := u.migrationService.DownOne()
	if err != nil {
		return errs.Wrap(errs.ServiceNA, "migration down failed", nil, err)
	}

	slog.Info("migrationService.DownOne success")
	return nil
}

func (u *Usecase) SetLogLevel(ctx context.Context, level string) error {
	if err := common.CheckPermission(ctx, []string{constant.RoleAdmin}); err != nil {
		return err
	}
	return u.loggingService.SetLevel(level)
}

func (u *Usecase) UpdateLogPolicy(ctx context.Context, method string, update logging.MethodLogPolicyUpdate) (logging.MethodLogPolicy, error) {
	if err := common.CheckPermission(ctx, []string{constant.RoleAdmin}); err != nil {
		return logging.MethodLogPolicy{}, err
	}
	return u.loggingService.UpdateMethodPolicy(method, update)
}

func (u *Usecase) DeleteLogPolicy(ctx context.Context, method string) error {
	if err := common.CheckPermission(ctx, []string{constant.RoleAdmin}); err != nil {
		return err
	}
	return u.loggingService.DeleteMethodPolicy(method)
}

func (u *Usecase) UpdateDefaultLogPolicy(ctx context.Context, update logging.MethodLogPolicyUpdate) (logging.MethodLogPolicy, error) {
	if err := common.CheckPermission(ctx, []string{constant.RoleAdmin}); err != nil {
		return logging.MethodLogPolicy{}, err
	}
	return u.loggingService.UpdateDefaultPolicy(update), nil
}
