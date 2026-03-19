package data

import (
	"context"

	auditCns "github.com/zhamspace/booking/internal/domain/audit/constants"
	"github.com/zhamspace/booking/internal/domain/common/dict/model"
)

type Repo struct {
}

func New() *Repo {
	return &Repo{}
}

func (r *Repo) Get(_ context.Context) *model.Main {
	result := &model.Main{
		Legend: make([]model.DictIdName, 0),
	}

	result.Legend = append(result.Legend, model.DictIdName{
		Id:   "audit_resources",
		Name: "Список источников аудит логов",
	})
	result.AuditResources = []model.DictIdName{
		{
			Id:   auditCns.ResourceSystem,
			Name: "Системные изменения",
		},
	}

	result.Legend = append(result.Legend, model.DictIdName{
		Id:   "audit_change_type",
		Name: "Вид измений данных / действие пользователя",
	})
	result.AuditChangeType = []model.DictIdName{
		{
			Id:   auditCns.ChangeTypeCreated,
			Name: "Создание",
		},
		{
			Id:   auditCns.ChangeTypeUpdate,
			Name: "Изменение",
		},
		{
			Id:   auditCns.ChangeTypeDelete,
			Name: "Удаление",
		},
	}

	return result
}
