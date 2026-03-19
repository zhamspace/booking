package model

import (
	"encoding/json"
	"time"

	commonModel "github.com/zhamspace/booking/internal/domain/common/model"
)

type Main struct {
	Id         string
	CreatedAt  time.Time
	User       *UserSt
	Resource   string
	ObjectId   string
	ChangeType string
	Object     json.RawMessage
}

type Edit struct {
	Id         *string
	CreatedAt  *time.Time
	UserId     *string
	ChangeType *string
	Resource   *string
	ObjectId   *string
	Object     *json.RawMessage
}

type ListReq struct {
	commonModel.ListParams

	Ids        *[]string
	UserId     *string
	Resource   *string
	ObjectId   *string
	ChangeType *string
}

type UserSt struct {
	Id        string
	Username  string
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time
}

type Job struct {
	UserId     string
	ChangeType string
	Resource   string
	ObjectId   string
	Object     any
}

type GetReq struct {
	Id string
}
