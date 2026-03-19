package model

import (
	"encoding/json"
	"time"

	"github.com/zhamspace/booking/internal/domain/audit/model"
)

type Select struct {
	Id         string
	CreatedAt  time.Time
	UserId     string
	Resource   string
	ObjectId   string
	ChangeType string
	Object     json.RawMessage
}

func (m *Select) ListColumnMap() map[string]any {
	return map[string]any{
		"id":          &m.Id,
		"created_at":  &m.CreatedAt,
		"user_id":     &m.UserId,
		"resource":    &m.Resource,
		"object_id":   &m.ObjectId,
		"change_type": &m.ChangeType,
		"object":      &m.Object,
	}
}

func (m *Select) PKColumnMap() map[string]any {
	return map[string]any{
		"id": &m.Id,
	}
}

func (m *Select) DefaultSortColumns() []string {
	return []string{
		"id",
	}
}

func EncodeMain(src *Select, _ int) *model.Main {
	if src == nil {
		return nil
	}
	return &model.Main{
		Id:        src.Id,
		CreatedAt: src.CreatedAt,
		User: &model.UserSt{
			Id: src.UserId,
		},
		Resource:   src.Resource,
		ObjectId:   src.ObjectId,
		ChangeType: src.ChangeType,
		Object:     src.Object,
	}
}
