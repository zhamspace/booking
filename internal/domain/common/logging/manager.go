package logging

import "github.com/zhamspace/booking/internal/errs"

type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) SetLevel(level string) error {
	if level == "" {
		return errs.New(errs.InvalidRequest, "log level is required", nil)
	}
	if err := SetLevel(level); err != nil {
		return errs.Wrap(errs.InvalidRequest, "invalid log level", map[string]string{
			"level": level,
		}, err)
	}
	return nil
}

func (m *Manager) UpdateMethodPolicy(method string, update MethodLogPolicyUpdate) (MethodLogPolicy, error) {
	if method == "" {
		return MethodLogPolicy{}, errs.New(errs.InvalidRequest, "method is required", nil)
	}
	return UpdateMethodPolicy(method, update), nil
}

func (m *Manager) DeleteMethodPolicy(method string) error {
	if method == "" {
		return errs.New(errs.InvalidRequest, "method is required", nil)
	}
	DeleteMethodPolicy(method)
	return nil
}

func (m *Manager) UpdateDefaultPolicy(update MethodLogPolicyUpdate) MethodLogPolicy {
	return UpdateDefaultPolicy(update)
}

func (m *Manager) GetDefaultPolicy() MethodLogPolicy {
	return GetDefaultPolicy()
}
