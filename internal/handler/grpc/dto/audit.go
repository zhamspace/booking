package dto

import (
	auditDomain "github.com/zhamspace/booking/internal/domain/audit"
	auditModel "github.com/zhamspace/booking/internal/domain/audit/model"

	"github.com/zhamspace/booking/pkg/proto/booking_v1"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuditEncoder struct {
	Serializer auditDomain.ObjectSerializer
}

func (e AuditEncoder) EncodeAuditMain(v *auditModel.Main, _ int) *booking_v1.AuditMain {
	if v == nil {
		return nil
	}

	res := &booking_v1.AuditMain{
		Id:         v.Id,
		CreatedAt:  timestamppb.New(v.CreatedAt),
		User:       EncodeUserSt(v.User, 0),
		Resource:   v.Resource,
		ObjectId:   v.ObjectId,
		ChangeType: v.ChangeType,
	}

	var decoded any
	if e.Serializer != nil && len(v.Object) > 0 {
		obj, err := e.Serializer.Deserialize(v.Resource, v.Object)
		if err == nil {
			decoded = obj
		}
	}

	switch v.Resource {
	}

	if res.Object == nil && decoded != nil {
		switch v := decoded.(type) {
		case map[string]any:
			if st, err := structpb.NewStruct(v); err == nil {
				res.Object = &booking_v1.AuditMain_Raw{Raw: st}
			}
		case map[string]string:
			raw := make(map[string]any, len(v))
			for key, value := range v {
				raw[key] = value
			}
			if st, err := structpb.NewStruct(raw); err == nil {
				res.Object = &booking_v1.AuditMain_Raw{Raw: st}
			}
		}
	}

	return res
}

func EncodeUserSt(v *auditModel.UserSt, _ int) *booking_v1.AccountMain {
	if v == nil {
		return nil
	}

	return &booking_v1.AccountMain{
		Id:        v.Id,
		Username:  v.Username,
		FirstName: v.FirstName,
		LastName:  v.LastName,
		Email:     v.Email,
		CreatedAt: timestamppb.New(v.CreatedAt),
	}
}
