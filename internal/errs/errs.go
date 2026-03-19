package errs

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zhamspace/booking/pkg/proto/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Err string

func (e Err) Error() string {
	return string(e)
}

// common errors
const (
	ServiceNA         = Err("service_not_available")
	NotImplemented    = Err("not_implemented")
	InvalidConfig     = Err("invalid_config")
	NoPermission      = Err("no_permission")
	ObjectNotFound    = Err("object_not_found")
	NoRows            = Err("err_no_rows")
	NotAuthorized     = Err("not_authorized")
	InvalidRequest    = Err("invalid_request")
	IncorrectPageSize = Err("incorrect_page_size")
)

type Error struct {
	Code    Err
	Message string
	Fields  map[string]string
	GRPC    codes.Code
	HTTP    int
	Cause   error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Code != "" {
		return e.Code.Error()
	}
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return "error"
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func New(code Err, message string, fields map[string]string) *Error {
	if fields == nil {
		fields = map[string]string{}
	}
	return &Error{
		Code:    code,
		Message: message,
		Fields:  fields,
	}
}

func Wrap(code Err, message string, fields map[string]string, cause error) *Error {
	err := New(code, message, fields)
	err.Cause = cause
	return err
}

func WithStatus(code Err, message string, fields map[string]string, grpcCode codes.Code, httpStatus int) *Error {
	err := New(code, message, fields)
	err.GRPC = grpcCode
	err.HTTP = httpStatus
	return err
}

func Validation(fields map[string]string) *Error {
	desc := "validation failed"
	if len(fields) > 0 {
		for k, v := range fields {
			desc = k + ": " + v
			if len(fields) > 1 {
				desc = fmt.Sprintf("%s (and %d more errors)", desc, len(fields)-1)
			}
			break
		}
	}

	return New(InvalidRequest, desc, fields)
}

func ToStatus(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := status.FromError(err); ok {
		return err
	}

	if appErr, ok := AsError(err); ok {
		stCode := grpcCode(appErr)
		detail := detailFromError(appErr)
		st := status.New(stCode, appErr.Message)
		stWithDetails, stErr := st.WithDetails(detail)
		if stErr == nil {
			return stWithDetails.Err()
		}
		return status.Error(stCode, appErr.Error())
	}

	if code, ok := asCode(err); ok {
		appErr := New(code, err.Error(), nil)
		stCode := grpcCode(appErr)
		detail := detailFromError(appErr)
		st := status.New(stCode, appErr.Message)
		stWithDetails, stErr := st.WithDetails(detail)
		if stErr == nil {
			return stWithDetails.Err()
		}
		return status.Error(stCode, appErr.Error())
	}

	return status.Error(codes.Internal, err.Error())
}

func ToHTTP(err error) (int, *common.ErrorDetail) {
	if err == nil {
		return http.StatusOK, nil
	}

	if appErr, ok := AsError(err); ok {
		return httpStatus(appErr), detailFromError(appErr)
	}

	if st, ok := status.FromError(err); ok {
		return httpStatusFromGRPC(st.Code()), detailFromStatus(st)
	}

	if code, ok := asCode(err); ok {
		appErr := New(code, err.Error(), nil)
		return httpStatus(appErr), detailFromError(appErr)
	}

	return http.StatusInternalServerError, &common.ErrorDetail{
		Code:    ServiceNA.Error(),
		Message: err.Error(),
	}
}

func AsError(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}
	if appErr, ok := errors.AsType[*Error](err); ok {
		return appErr, true
	}
	return nil, false
}

func ExtractDetail(err error) (*common.ErrorDetail, codes.Code) {
	if err == nil {
		return &common.ErrorDetail{
			Code:    "",
			Message: "",
		}, codes.OK
	}

	if appErr, ok := AsError(err); ok {
		return detailFromError(appErr), grpcCode(appErr)
	}

	if st, ok := status.FromError(err); ok {
		return detailFromStatus(st), st.Code()
	}

	if code, ok := asCode(err); ok {
		appErr := New(code, err.Error(), nil)
		return detailFromError(appErr), grpcCode(appErr)
	}

	return &common.ErrorDetail{
		Code:    ServiceNA.Error(),
		Message: err.Error(),
	}, codes.Internal
}

func asCode(err error) (Err, bool) {
	if code, ok := errors.AsType[Err](err); ok {
		return code, true
	}
	return "", false
}

func detailFromError(err *Error) *common.ErrorDetail {
	if err == nil {
		return &common.ErrorDetail{}
	}
	code := err.Code.Error()
	if err.Code == "" {
		code = ServiceNA.Error()
	}
	return &common.ErrorDetail{
		Code:    code,
		Message: err.Message,
		Fields:  err.Fields,
	}
}

func detailFromStatus(st *status.Status) *common.ErrorDetail {
	if st == nil {
		return &common.ErrorDetail{Code: ServiceNA.Error(), Message: "unknown error"}
	}
	if len(st.Details()) > 0 {
		if detail, ok := st.Details()[0].(*common.ErrorDetail); ok {
			return detail
		}
	}
	return &common.ErrorDetail{
		Code:    ServiceNA.Error(),
		Message: st.Message(),
	}
}

func grpcCode(err *Error) codes.Code {
	if err == nil {
		return codes.Internal
	}
	if err.GRPC != 0 {
		return err.GRPC
	}
	return errToStatusCode(err.Code)
}

func httpStatus(err *Error) int {
	if err == nil {
		return http.StatusInternalServerError
	}
	if err.HTTP != 0 {
		return err.HTTP
	}
	return httpStatusFromErr(err.Code)
}

func errToStatusCode(err error) codes.Code {
	switch {
	case errors.Is(err, ServiceNA):
		return codes.Unavailable
	case errors.Is(err, NotImplemented):
		return codes.Unimplemented
	case errors.Is(err, NoRows), errors.Is(err, ObjectNotFound):
		return codes.NotFound
	case errors.Is(err, NoPermission):
		return codes.PermissionDenied
	case errors.Is(err, NotAuthorized):
		return codes.Unauthenticated
	case errors.Is(err, InvalidConfig), errors.Is(err, IncorrectPageSize), errors.Is(err, InvalidRequest):
		return codes.InvalidArgument
	default:
		return codes.InvalidArgument
	}
}

func httpStatusFromErr(err error) int {
	switch {
	case errors.Is(err, ServiceNA):
		return http.StatusInternalServerError
	case errors.Is(err, NotImplemented):
		return http.StatusNotImplemented
	case errors.Is(err, NoRows), errors.Is(err, ObjectNotFound):
		return http.StatusNotFound
	case errors.Is(err, NoPermission):
		return http.StatusForbidden
	case errors.Is(err, NotAuthorized):
		return http.StatusUnauthorized
	case errors.Is(err, InvalidConfig), errors.Is(err, IncorrectPageSize), errors.Is(err, InvalidRequest):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func httpStatusFromGRPC(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return 499
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
