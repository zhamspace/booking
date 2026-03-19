package dto

import (
	dictModel "github.com/zhamspace/booking/internal/domain/common/dict/model"
	"github.com/zhamspace/booking/pkg/proto/booking_v1"
)

func EncodeDictMain(v *dictModel.Main) *booking_v1.DictRep {
	if v == nil {
		return nil
	}

	result := &booking_v1.DictRep{}

	result.Legend = encodeDictIdNames(v.Legend)
	result.AuditResources = encodeDictIdNames(v.AuditResources)
	result.AuditChangeType = encodeDictIdNames(v.AuditChangeType)

	return result
}

func EncodeDictIdName(v *dictModel.DictIdName) *booking_v1.DictIdName {
	if v == nil {
		return nil
	}

	return &booking_v1.DictIdName{
		Id:   v.Id,
		Name: &v.Name,
	}
}

func encodeDictIdNames(items []dictModel.DictIdName) []*booking_v1.DictIdName {
	if items == nil {
		return nil
	}

	result := make([]*booking_v1.DictIdName, len(items))
	for i, item := range items {
		result[i] = EncodeDictIdName(&item)
	}
	return result
}
