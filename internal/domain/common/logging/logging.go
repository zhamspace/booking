package logging

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
)

type MethodLogPolicy struct {
	LogSuccess  bool `json:"log_success"`
	LogError    bool `json:"log_error"`
	LogRequest  bool `json:"log_request"`
	LogResponse bool `json:"log_response"`
}

type MethodLogPolicyUpdate struct {
	LogSuccess  *bool `json:"log_success,omitempty"`
	LogError    *bool `json:"log_error,omitempty"`
	LogRequest  *bool `json:"log_request,omitempty"`
	LogResponse *bool `json:"log_response,omitempty"`
}

var (
	levelVar       slog.LevelVar
	mu             sync.RWMutex
	defaultPolicy  = MethodLogPolicy{LogError: true}
	methodPolicies = map[string]MethodLogPolicy{}
)

func LevelVar() *slog.LevelVar {
	return &levelVar
}

func InitLevel(level string) error {
	parsed, err := parseLevel(level)
	if err != nil {
		return err
	}
	levelVar.Set(parsed)
	return nil
}

func SetLevel(level string) error {
	parsed, err := parseLevel(level)
	if err != nil {
		return err
	}
	levelVar.Set(parsed)
	return nil
}

func GetLevel() slog.Level {
	return levelVar.Level()
}

func GetPolicy(method string) MethodLogPolicy {
	mu.RLock()
	defer mu.RUnlock()
	policy := defaultPolicy
	if methodPolicy, ok := methodPolicies[method]; ok {
		policy = methodPolicy
	}
	return policy
}

func UpdateMethodPolicy(method string, update MethodLogPolicyUpdate) MethodLogPolicy {
	mu.Lock()
	defer mu.Unlock()

	policy := defaultPolicy
	if existing, ok := methodPolicies[method]; ok {
		policy = existing
	}

	applyUpdate(&policy, update)
	methodPolicies[method] = policy
	return policy
}

func DeleteMethodPolicy(method string) {
	mu.Lock()
	defer mu.Unlock()
	delete(methodPolicies, method)
}

func UpdateDefaultPolicy(update MethodLogPolicyUpdate) MethodLogPolicy {
	mu.Lock()
	defer mu.Unlock()
	applyUpdate(&defaultPolicy, update)
	return defaultPolicy
}

func GetDefaultPolicy() MethodLogPolicy {
	mu.RLock()
	defer mu.RUnlock()
	return defaultPolicy
}

func applyUpdate(policy *MethodLogPolicy, update MethodLogPolicyUpdate) {
	if update.LogSuccess != nil {
		policy.LogSuccess = *update.LogSuccess
	}
	if update.LogError != nil {
		policy.LogError = *update.LogError
	}
	if update.LogRequest != nil {
		policy.LogRequest = *update.LogRequest
	}
	if update.LogResponse != nil {
		policy.LogResponse = *update.LogResponse
	}
}

func parseLevel(level string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info", "":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("invalid log level: %s", level)
	}
}
