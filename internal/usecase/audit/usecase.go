package audit

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/zhamspace/booking/internal/constant"
	"github.com/zhamspace/booking/internal/domain/common/tools"
	"github.com/zhamspace/booking/internal/errs"

	auditCns "github.com/zhamspace/booking/internal/domain/audit/constants"
	auditModel "github.com/zhamspace/booking/internal/domain/audit/model"
	accountModel "github.com/zhamspace/booking/internal/service/account/model"

	"github.com/samber/lo"
)

type Usecase struct {
	accountService AccountServiceI
	auditService   ServiceI
}

func New(
	accountService AccountServiceI,
	auditService ServiceI,
) *Usecase {
	return &Usecase{
		accountService: accountService,
		auditService:   auditService,
	}
}

func (u *Usecase) List(ctx context.Context, pars *auditModel.ListReq) (_ []*auditModel.Main, _ int64, finalError error) {
	if pars == nil {
		return nil, 0, errs.New(errs.InvalidConfig, "list params are required", nil)
	}

	fields := map[string]string{}
	if !pars.OnlyCount {
		if err := tools.RequirePageSize(pars.ListParams, constant.MaxPageSize); err != nil {
			fields["list_params.page_size"] = fmt.Sprintf("must be between 1 and %d", constant.MaxPageSize)
		}
	}

	if pars.Resource != nil {
		validValues := map[string]struct{}{
			auditCns.ResourceSystem:  {},
			auditCns.ResourceBooking: {},
		}
		if _, ok := validValues[*pars.Resource]; !ok {
			fields["resource"] = fmt.Sprintf("invalid resource: %s", *pars.Resource)
		}
	}

	if len(fields) > 0 {
		return nil, 0, errs.Validation(fields)
	}

	items, totalCount, err := u.auditService.List(ctx, pars)
	if err != nil {
		fields["method"] = "auditService.List"
		return nil, 0, errs.Wrap(errs.ServiceNA, "failed data retrieval", fields, err)
	}

	userMap := make(map[string]*accountModel.Main)
	userIds := make([]string, 0, len(items))
	userSet := map[string]struct{}{}
	lo.ForEach(items, func(item *auditModel.Main, _ int) {
		if item == nil {
			return
		}
		if item.User != nil {
			if _, ok := userSet[item.User.Id]; !ok {
				userSet[item.User.Id] = struct{}{}
				userIds = append(userIds, item.User.Id)
			}
		}
	})

	if len(userIds) > 0 {
		for _, userId := range userIds {
			user, found, err := u.accountService.GetUser(ctx, &accountModel.GetReq{Id: userId}, false)
			if err != nil {
				slog.ErrorContext(ctx, "failed get user", "userId", userId, "error", err)
			}
			if found {
				userMap[userId] = user
			}
		}
	}

	lo.ForEach(items, func(item *auditModel.Main, _ int) {
		if item == nil {
			return
		}
		if item.User != nil {
			if user, ok := userMap[item.User.Id]; ok {
				item.User.Username = user.Username
				item.User.FirstName = user.FirstName
				item.User.LastName = user.LastName
				item.User.Email = user.Email
				item.User.CreatedAt = user.CreatedAt
			}
		}
	})

	return items, totalCount, nil
}

func (u *Usecase) Get(ctx context.Context, pars *auditModel.GetReq) (*auditModel.Main, error) {
	fields := map[string]string{}

	item, _, err := u.auditService.Get(ctx, pars, true)
	if err != nil {
		fields["method"] = "auditService.Get"
		return nil, errs.Wrap(errs.ServiceNA, "failed data retrieval", fields, err)
	}

	if item.User != nil {
		user, found, err := u.accountService.GetUser(ctx, &accountModel.GetReq{Id: item.User.Id}, false)
		if err != nil {
			slog.ErrorContext(ctx, "failed get user", "userId", item.User.Id, "error", err)
		}
		if found {
			item.User.Username = user.Username
			item.User.FirstName = user.FirstName
			item.User.LastName = user.LastName
			item.User.Email = user.Email
			item.User.CreatedAt = user.CreatedAt
		}
	}

	return item, nil
}
