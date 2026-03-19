package db

import "github.com/zhamspace/booking/internal/domain/audit/model"

var allowedSortFields = map[string]string{
	"id":         "id",
	"created_at": "created_at",
}

func (r *Repo) getConditions(pars *model.ListReq) (map[string]any, map[string][]any) {
	conditions := make(map[string]any)
	conditionExps := make(map[string][]any)

	if pars.Ids != nil {
		conditions["id"] = *pars.Ids
	}
	if pars.ObjectId != nil {
		conditions["object_id"] = *pars.ObjectId
	}
	if pars.Resource != nil {
		conditions["resource"] = *pars.Resource
	}

	return conditions, conditionExps
}
