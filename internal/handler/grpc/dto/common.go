package dto

import (
	commonModel "github.com/zhamspace/booking/internal/domain/common/model"

	"github.com/zhamspace/booking/pkg/proto/common"
)

func DecodeListParams(listParams *common.ListParamsSt) commonModel.ListParams {
	if listParams == nil {
		return commonModel.ListParams{}
	}

	return commonModel.ListParams{
		Page:           listParams.Page,
		PageSize:       listParams.PageSize,
		WithTotalCount: listParams.WithTotalCount,
		OnlyCount:      listParams.OnlyCount,
		SortName:       listParams.SortName,
		Sort:           listParams.Sort,
	}
}

func EncodePaginationInfoSt(listParams *common.ListParamsSt, totalCount int64) *common.PaginationInfoSt {
	if listParams == nil {
		return &common.PaginationInfoSt{}
	}

	return &common.PaginationInfoSt{
		Page:       listParams.Page,
		PageSize:   listParams.PageSize,
		TotalCount: totalCount,
	}
}

func ToSlicePtr[T any](v []T) *[]T {
	if v == nil {
		return nil
	}

	return &v
}
