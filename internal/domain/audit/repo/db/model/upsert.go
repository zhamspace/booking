package model

import (
	"encoding/json"
	"time"

	"github.com/zhamspace/booking/internal/domain/audit/model"
)

type Upsert struct {
	PKId string

	Id         *string
	CreatedAt  *time.Time
	UserId     *string
	Resource   *string
	ObjectId   *string
	ChangeType *string
	Object     *json.RawMessage
}

func (m *Upsert) CreateColumnMap() map[string]any {
	result := map[string]any{}

	if m.Id != nil {
		result["id"] = *m.Id
	}

	if m.CreatedAt != nil {
		result["created_at"] = *m.CreatedAt
	}

	if m.UserId != nil {
		result["user_id"] = *m.UserId
	}

	if m.Resource != nil {
		result["resource"] = *m.Resource
	}

	if m.ObjectId != nil {
		result["object_id"] = *m.ObjectId
	}

	if m.ChangeType != nil {
		result["change_type"] = *m.ChangeType
	}

	if m.Object != nil {
		result["object"] = *m.Object
	}

	return result
}

func (m *Upsert) ReturningColumnMap() map[string]any {
	return map[string]any{
		"id": &m.PKId,
	}
}

func (m *Upsert) UpdateColumnMap() map[string]any {
	return m.CreateColumnMap()
}

func (m *Upsert) PKColumnMap() map[string]any {
	return map[string]any{
		"id": m.PKId,
	}
}

func DecodeEdit(src *model.Edit) *Upsert {
	if src == nil {
		return &Upsert{}
	}

	return &Upsert{
		Id:         src.Id,
		CreatedAt:  src.CreatedAt,
		UserId:     src.UserId,
		ChangeType: src.ChangeType,
		Resource:   src.Resource,
		ObjectId:   src.ObjectId,
		Object:     src.Object,
	}
}
