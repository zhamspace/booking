package http

import (
	"encoding/json"
	"net/http"

	"github.com/zhamspace/booking/internal/errs"
	"github.com/zhamspace/booking/pkg/proto/common"
)

// errorResponse sends a normalized error envelope that matches gRPC ErrorDetail.
func errorResponse(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	finalStatus, detail := errs.ToHTTP(err)
	if status != 0 {
		finalStatus = status
	}

	payload := detail
	if payload == nil {
		payload = &common.ErrorDetail{
			Code:    errs.ServiceNA.Error(),
			Message: err.Error(),
		}
	}

	result, _ := json.Marshal(payload)
	w.WriteHeader(finalStatus)
	_, _ = w.Write(result)
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return errs.New(errs.InvalidRequest, "invalid json body", map[string]string{
			"error": err.Error(),
		})
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	if data, err := json.Marshal(payload); err == nil {
		_, _ = w.Write(data)
	}
}
