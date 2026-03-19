package tools

import (
	"slices"

	"github.com/zhamspace/booking/internal/domain/common/model"
	"github.com/zhamspace/booking/internal/errs"
)

const (
	defaultMaxPageSize int64 = 100
)

func RequirePageSize(pars model.ListParams, maxPageSize int64) error {
	if maxPageSize == 0 {
		maxPageSize = defaultMaxPageSize
	}

	if pars.PageSize <= 0 || pars.PageSize > maxPageSize {
		return errs.IncorrectPageSize
	}

	return nil
}

func CompareLists[T1 any, T2 any](oldList, newList []T1, getId func(T1) string, getChanges func(oldItem T1, newItem T1) (T2, bool)) (deleteItems map[string]T1, updateItems map[string]T2, insertItems map[string]T1) {
	dList := make(map[string]T1, len(oldList))
	uList := make(map[string]T2, len(oldList))
	iList := make(map[string]T1, len(newList))

	oldMap := map[string]T1{}
	for _, oldItem := range oldList {
		oldMap[getId(oldItem)] = oldItem
	}

	newMap := map[string]T1{}
	for _, newItem := range newList {
		newMap[getId(newItem)] = newItem
	}

	var changes T2
	var oldItem T1
	var newItem T1
	var id string
	var ok bool

	for id, oldItem = range oldMap {
		if newItem, ok = newMap[id]; ok {
			if changes, ok = getChanges(oldItem, newItem); ok {
				uList[id] = changes
			}
		} else {
			dList[id] = oldItem
		}
	}

	for id, newItem = range newMap {
		if _, ok = oldMap[id]; !ok {
			iList[id] = newItem
		}
	}

	return dList, uList, iList
}

func SliceHasValue[T comparable](sl []T, v T) bool {
	return slices.Contains(sl, v)
}

func SliceHasValueBy[T comparable](sl []T, v T, cmp func(T, T) bool) bool {
	for _, x := range sl {
		if cmp(x, v) {
			return true
		}
	}

	return false
}

func SlicesAreSame[T comparable](a, b []T) bool {
	for _, x := range a {
		if !SliceHasValue(b, x) {
			return false
		}
	}

	for _, x := range b {
		if !SliceHasValue(a, x) {
			return false
		}
	}

	return true
}

func SlicesAreSameBy[T comparable](a, b []T, cmp func(T, T) bool) bool {
	for _, x := range a {
		if !SliceHasValueBy(b, x, cmp) {
			return false
		}
	}

	for _, x := range b {
		if !SliceHasValueBy(a, x, cmp) {
			return false
		}
	}

	return true
}
