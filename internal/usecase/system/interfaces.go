package system

import "github.com/zhamspace/booking/internal/domain/common/logging"

type LoggingServiceI interface {
	SetLevel(level string) error
	UpdateMethodPolicy(method string, update logging.MethodLogPolicyUpdate) (logging.MethodLogPolicy, error)
	DeleteMethodPolicy(method string) error
	UpdateDefaultPolicy(update logging.MethodLogPolicyUpdate) logging.MethodLogPolicy
	GetDefaultPolicy() logging.MethodLogPolicy
}
