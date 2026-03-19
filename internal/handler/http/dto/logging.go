package dto

type LogLevelReq struct {
	Level string `json:"level"`
}

type LogPolicyReq struct {
	Method      string `json:"method"`
	LogSuccess  *bool  `json:"log_success,omitempty"`
	LogError    *bool  `json:"log_error,omitempty"`
	LogRequest  *bool  `json:"log_request,omitempty"`
	LogResponse *bool  `json:"log_response,omitempty"`
}

type LogPolicyUpdateReq struct {
	LogSuccess  *bool `json:"log_success,omitempty"`
	LogError    *bool `json:"log_error,omitempty"`
	LogRequest  *bool `json:"log_request,omitempty"`
	LogResponse *bool `json:"log_response,omitempty"`
}
