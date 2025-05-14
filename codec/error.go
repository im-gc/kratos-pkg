package codec

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/im-gc/kratos-pkg/logger"

	"github.com/go-kratos/kratos/v2/errors"
)

var (
	ErrUnauthorized = errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	ErrURINotFound  = errors.Unauthorized("NOT_FOUND", "URI not found")
	ErrTokenInvalid = errors.Unauthorized("INVALID", "Token invalid")
)

type ErrInternalServer struct {
	e *errors.Error
}

func (e ErrInternalServer) Error() string {
	return e.e.Error()
}

func NewErrInternalServer(message string) ErrInternalServer {
	return ErrInternalServer{
		e: errors.New(500, "INTERNAL_SERVER_ERROR", message),
	}
}

// DefaultErrorEncoder encodes the error to the HTTP response.
func DefaultErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {

	resp := &Response{
		Code:    -1,
		Data:    "",
		Message: err.Error(),
	}
	body, e := json.Marshal(resp)
	if nil != e {
		logger.ErrorfWithContext(context.Background(), "pkg.codec: %s", e.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	statusCode := http.StatusBadRequest
	if errors.Is(err, ErrTokenInvalid) || errors.Is(err, ErrURINotFound) || errors.Is(err, ErrUnauthorized) {
		statusCode = http.StatusUnauthorized
	}
	switch err.(type) {
	case ErrInternalServer:
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)
	_, _ = w.Write(body)
}
