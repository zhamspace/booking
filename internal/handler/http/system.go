package http

import (
	"net/http"
	"strconv"

	"github.com/zhamspace/booking/internal/handler/http/dto"
	usecaseSystemP "github.com/zhamspace/booking/internal/usecase/system"
)

type System struct {
	usecase *usecaseSystemP.Usecase
}

func NewSystem(usecase *usecaseSystemP.Usecase) *System {
	return &System{
		usecase: usecase,
	}
}

func (h *System) Command(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	cmdNum, err := strconv.Atoi(pathParams["num"])
	if err != nil {
		errorResponse(w, 0, err)
		return
	}

	err = h.usecase.Command(r.Context(), cmdNum)
	if err != nil {
		errorResponse(w, 0, err)
		return
	}
}

func (h *System) MigrationUp(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	err := h.usecase.MigrationUp(r.Context())
	if err != nil {
		errorResponse(w, 0, err)
		return
	}
}

func (h *System) MigrationDownOne(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	err := h.usecase.MigrationDownOne(r.Context())
	if err != nil {
		errorResponse(w, 0, err)
		return
	}
}

func (h *System) SetLogLevel(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var req dto.LogLevelReq
	if err := decodeJSON(r, &req); err != nil {
		errorResponse(w, 0, err)
		return
	}

	if err := h.usecase.SetLogLevel(r.Context(), req.Level); err != nil {
		errorResponse(w, 0, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"level": req.Level,
	})
}

func (h *System) UpdateLogPolicy(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var req dto.LogPolicyReq
	if err := decodeJSON(r, &req); err != nil {
		errorResponse(w, 0, err)
		return
	}

	policy, err := h.usecase.UpdateLogPolicy(r.Context(), req.Method, usecaseSystemP.MethodLogPolicyUpdate{
		LogSuccess:  req.LogSuccess,
		LogError:    req.LogError,
		LogRequest:  req.LogRequest,
		LogResponse: req.LogResponse,
	})
	if err != nil {
		errorResponse(w, 0, err)
		return
	}

	writeJSON(w, http.StatusOK, policy)
}

func (h *System) DeleteLogPolicy(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	if err := h.usecase.DeleteLogPolicy(r.Context(), pathParams["method"]); err != nil {
		errorResponse(w, 0, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *System) UpdateDefaultLogPolicy(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	var req dto.LogPolicyUpdateReq
	if err := decodeJSON(r, &req); err != nil {
		errorResponse(w, 0, err)
		return
	}

	policy, err := h.usecase.UpdateDefaultLogPolicy(r.Context(), usecaseSystemP.MethodLogPolicyUpdate{
		LogSuccess:  req.LogSuccess,
		LogError:    req.LogError,
		LogRequest:  req.LogRequest,
		LogResponse: req.LogResponse,
	})
	if err != nil {
		errorResponse(w, 0, err)
		return
	}

	writeJSON(w, http.StatusOK, policy)
}
