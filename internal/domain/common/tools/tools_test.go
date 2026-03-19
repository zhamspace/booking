package tools

import (
	"testing"

	"github.com/zhamspace/booking/internal/constant"
	"github.com/zhamspace/booking/internal/domain/common/model"
	"github.com/zhamspace/booking/internal/errs"

	"github.com/stretchr/testify/require"
)

func TestRequirePageSize(t *testing.T) {
	cases := []struct {
		name        string
		pars        model.ListParams
		maxPageSize int64
		want        error
	}{
		{
			name:        "1",
			pars:        model.ListParams{PageSize: 10},
			maxPageSize: constant.MaxPageSize,
			want:        nil,
		},
		{
			name:        "2",
			pars:        model.ListParams{PageSize: 5},
			maxPageSize: constant.MaxPageSize,
			want:        nil,
		},
		{
			name:        "3",
			pars:        model.ListParams{PageSize: 10},
			maxPageSize: 0, // defaultMaxPageSize = 100
			want:        nil,
		},
		{
			name:        "4",
			pars:        model.ListParams{PageSize: defaultMaxPageSize + 1},
			maxPageSize: 0, // defaultMaxPageSize = 100
			want:        errs.IncorrectPageSize,
		},
		{
			name:        "5",
			pars:        model.ListParams{PageSize: 0},
			maxPageSize: constant.MaxPageSize,
			want:        errs.IncorrectPageSize,
		},
		{
			name:        "6",
			pars:        model.ListParams{PageSize: constant.MaxPageSize + 1},
			maxPageSize: constant.MaxPageSize,
			want:        errs.IncorrectPageSize,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := RequirePageSize(c.pars, c.maxPageSize)
			require.Equal(t, c.want, got)
		})
	}
}

func TestCompareLists(t *testing.T) {
	type itemType struct {
		Id   string
		Name string
	}

	type itemChangeType struct {
		Id   *string
		Name *string
	}

	getId := func(item *itemType) string {
		return item.Id
	}

	getChanges := func(oldItem *itemType, newItem *itemType) (*itemChangeType, bool) {
		changes := &itemChangeType{}
		if oldItem.Id != newItem.Id {
			changes.Id = &newItem.Id
		}
		if oldItem.Name != newItem.Name {
			changes.Name = &newItem.Name
		}
		if *changes != (itemChangeType{}) {
			return changes, true
		}
		return nil, false
	}

	cases := []struct {
		name           string
		oldList        []*itemType
		newList        []*itemType
		wantDeleteList map[string]*itemType
		wantUpdateList map[string]*itemChangeType
		wantInsertList map[string]*itemType
	}{
		{
			name: "1",
			oldList: []*itemType{
				{Id: "1", Name: "1"},
				{Id: "2", Name: "2"},
				{Id: "3", Name: "3"},
			},
			newList: []*itemType{
				{Id: "1", Name: "1"},
				{Id: "2", Name: "2"},
				{Id: "4", Name: "4"},
			},
			wantDeleteList: map[string]*itemType{
				"3": {Id: "3", Name: "3"},
			},
			wantUpdateList: map[string]*itemChangeType{},
			wantInsertList: map[string]*itemType{
				"4": {Id: "4", Name: "4"},
			},
		},
		{
			name: "2",
			oldList: []*itemType{
				{Id: "1", Name: "1"},
				{Id: "2", Name: "2"},
				{Id: "3", Name: "3"},
			},
			newList: []*itemType{
				{Id: "1", Name: "1"},
				{Id: "2", Name: "2"},
				{Id: "3", Name: "3"},
			},
			wantDeleteList: map[string]*itemType{},
			wantUpdateList: map[string]*itemChangeType{},
			wantInsertList: map[string]*itemType{},
		},
		{
			name: "3",
			oldList: []*itemType{
				{Id: "1", Name: "1"},
				{Id: "2", Name: "2"},
				{Id: "3", Name: "3"},
			},
			newList: []*itemType{
				{Id: "1", Name: "1"},
				{Id: "3", Name: "4"},
				{Id: "5", Name: "5"},
			},
			wantDeleteList: map[string]*itemType{
				"2": {Id: "2", Name: "2"},
			},
			wantUpdateList: map[string]*itemChangeType{
				"3": {Name: new("4")},
			},
			wantInsertList: map[string]*itemType{
				"5": {Id: "5", Name: "5"},
			},
		},
		{
			name: "4",
			oldList: []*itemType{
				{Id: "1", Name: "1"},
				{Id: "2", Name: "2"},
			},
			newList: []*itemType{},
			wantDeleteList: map[string]*itemType{
				"1": {Id: "1", Name: "1"},
				"2": {Id: "2", Name: "2"},
			},
			wantUpdateList: map[string]*itemChangeType{},
			wantInsertList: map[string]*itemType{},
		},
		{
			name:    "5",
			oldList: []*itemType{},
			newList: []*itemType{
				{Id: "1", Name: "1"},
				{Id: "2", Name: "2"},
			},
			wantDeleteList: map[string]*itemType{},
			wantUpdateList: map[string]*itemChangeType{},
			wantInsertList: map[string]*itemType{
				"1": {Id: "1", Name: "1"},
				"2": {Id: "2", Name: "2"},
			},
		},
		{
			name: "6",
			oldList: []*itemType{
				{Id: "1", Name: "1"},
				{Id: "2", Name: "2"},
			},
			newList: []*itemType{
				{Id: "1", Name: "1c"},
				{Id: "2", Name: "2c"},
			},
			wantDeleteList: map[string]*itemType{},
			wantUpdateList: map[string]*itemChangeType{
				"1": {Name: new("1c")},
				"2": {Name: new("2c")},
			},
			wantInsertList: map[string]*itemType{},
		},
		{
			name:           "7",
			oldList:        []*itemType{},
			newList:        []*itemType{},
			wantDeleteList: map[string]*itemType{},
			wantUpdateList: map[string]*itemChangeType{},
			wantInsertList: map[string]*itemType{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotDeleteList, gotUpdateList, gotInsertList := CompareLists(c.oldList, c.newList, getId, getChanges)
			require.Equal(t, c.wantDeleteList, gotDeleteList)
			require.Equal(t, c.wantUpdateList, gotUpdateList)
			require.Equal(t, c.wantInsertList, gotInsertList)
		})
	}
}
