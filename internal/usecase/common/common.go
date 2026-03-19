package common

import (
	"context"
	"strings"

	"github.com/zhamspace/booking/internal/domain/common/session"
	"github.com/zhamspace/booking/internal/errs"
)

func CheckPermission(ctx context.Context, validRoles []string) error {
	if len(validRoles) == 0 {
		return nil
	}

	s := session.ExtractFromContext(ctx)
	if hasRole(s.Roles, validRoles) {
		return nil
	}

	if s.Sub == "" {
		return errs.New(errs.NotAuthorized, "invalid session", nil)
	}

	if len(s.Roles) == 0 {
		return errs.New(errs.NotAuthorized, "roles not found", map[string]string{
			"sub":      s.Sub,
			"required": strings.Join(validRoles, ","),
		})
	}
	return errs.New(errs.NoPermission, "required role not found in session", map[string]string{
		"required":  strings.Join(validRoles, ","),
		"presented": strings.Join(s.Roles, ","),
	})
}

func hasRole(roles, validRoles []string) bool {
	if len(roles) == 0 {
		return false
	}

	valid := make(map[string]struct{}, len(validRoles))
	for _, role := range validRoles {
		valid[role] = struct{}{}
	}

	for _, role := range roles {
		if _, ok := valid[role]; ok {
			return true
		}
	}

	return false
}
