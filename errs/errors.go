package errs

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type Status string

const (
	INTERNAL_SERVER_ERROR Status = "INTERNAL_SERVER_ERROR"
	NOT_FOUND             Status = "NOT_FOUND"
	BAD_REQUEST           Status = "BAD_REQUEST"
	UNPROCESSABLE_ENTITY  Status = "UNPROCESSABLE_ENTITY"
	CONFLICT              Status = "CONFLICT"
	UNAUTHORIZED          Status = "UNAUTHORIZED"
	SERVICE_UNAVAILABLE   Status = "SERVICE_UNAVAILABLE"
	NO_CONTENT            Status = "NO_CONTENT"
)

type ErrorCode string

const (
	SystemError           ErrorCode = "internal_error"
	BadRequestBody        ErrorCode = "bad_request_body"
	MissingRequiredFields ErrorCode = "missing_req_fields"
	BadRequest            ErrorCode = "bad_request"
	NoContent             ErrorCode = "no_content"
	UnprocessableEntity   ErrorCode = "unprocessable_entity"
	NotFound              ErrorCode = "not_found"
	ServiceUnavailable    ErrorCode = "service_unavailable"
)

// ErrorService is a custom error type
type ErrorService struct {
	StatusText   string
	StatusCode   int
	ErrorCode    string
	ErrorMessage string
	StackTrace   string // Store stack trace as a string
}

// Error implements the error interface for ErrorService
func (e *ErrorService) Error() string {
	return fmt.Sprintf(" Error: %s", e.ErrorMessage)
}

// GetStackTrace returns the stored stack trace
func (e *ErrorService) GetStackTrace() string {
	return formatStackTrace(e.StackTrace) // Format errors.StackTrace for pretty printing
}

// captureStackTrace captures the stack trace using github.com/pkg/errors
func captureStackTrace(message string) string {
	return fmt.Sprintf("%+v", errors.WithStack(errors.New(message)))
}

func formatStackTrace(trace string) string {
	lines := strings.Split(trace, "\n")
	var formattedLines []string
	for _, line := range lines {
		formattedLines = append(formattedLines, strings.TrimSpace(line))
	}
	return strings.Join(formattedLines, "\n")
}

func GetStatusCode(errCode Status) int {
	switch errCode {
	case INTERNAL_SERVER_ERROR:
		return http.StatusInternalServerError
	case NOT_FOUND:
		return http.StatusNotFound
	case CONFLICT:
		return http.StatusConflict
	case UNAUTHORIZED:
		return http.StatusUnauthorized
	case BAD_REQUEST:
		return http.StatusBadRequest
	case UNPROCESSABLE_ENTITY:
		return http.StatusUnprocessableEntity
	case SERVICE_UNAVAILABLE:
		return http.StatusServiceUnavailable
	case NO_CONTENT:
		return http.StatusNoContent
	default:
		return http.StatusInternalServerError
	}
}

// CreateError creates a new ErrorService instance with the given code and message
func CreateError(status Status, code ErrorCode, message string) *ErrorService {
	statusCode := GetStatusCode(status)
	stackTrace := captureStackTrace(message)
	return &ErrorService{
		StatusText:   http.StatusText(statusCode),
		StatusCode:   statusCode,
		ErrorCode:    string(code),
		ErrorMessage: message,
		StackTrace:   stackTrace,
	}
}

// CreateInternalError creates a new ErrorService instance for internal server errors
func CreateInternalError() *ErrorService {
	statusCode := GetStatusCode(INTERNAL_SERVER_ERROR)
	stackTrace := captureStackTrace("Internal Server Error")
	return &ErrorService{
		StatusText:   http.StatusText(statusCode),
		StatusCode:   statusCode,
		ErrorCode:    string(SystemError),
		ErrorMessage: "Internal Server Error",
		StackTrace:   stackTrace,
	}
}

// CreateServiceUnavailableError creates a new ErrorService instance for service unavailable errors
func CreateServiceUnavailableError(message string) *ErrorService {
	statusCode := GetStatusCode(SERVICE_UNAVAILABLE)
	stackTrace := captureStackTrace(message)
	return &ErrorService{
		StatusText:   http.StatusText(statusCode),
		StatusCode:   statusCode,
		ErrorCode:    string(SystemError),
		ErrorMessage: message,
		StackTrace:   stackTrace,
	}
}

// CreateBadRequestError creates a new ErrorService instance for bad request errors
func CreateBadRequestError(code ErrorCode, message string) *ErrorService {
	statusCode := GetStatusCode(BAD_REQUEST)
	stackTrace := captureStackTrace(message)
	return &ErrorService{
		StatusText:   http.StatusText(statusCode),
		StatusCode:   statusCode,
		ErrorCode:    string(code),
		ErrorMessage: message,
		StackTrace:   stackTrace,
	}
}

// CreateNotFoundError creates a new ErrorService instance for not found errors
func CreateNotFoundError(code ErrorCode, message string) *ErrorService {
	statusCode := GetStatusCode(NOT_FOUND)
	stackTrace := captureStackTrace(message)
	return &ErrorService{
		StatusText:   http.StatusText(statusCode),
		StatusCode:   statusCode,
		ErrorCode:    string(code),
		ErrorMessage: message,
		StackTrace:   stackTrace,
	}
}

// CreateUnprocessableEntityError creates a new ErrorService instance for unprocessable entity errors
func CreateUnprocessableEntityError(code ErrorCode, message string) *ErrorService {
	statusCode := GetStatusCode(UNPROCESSABLE_ENTITY)
	stackTrace := captureStackTrace(message)
	return &ErrorService{
		StatusText:   http.StatusText(statusCode),
		StatusCode:   statusCode,
		ErrorCode:    string(code),
		ErrorMessage: message,
		StackTrace:   stackTrace,
	}
}

func CreateNoContentError(code ErrorCode, message string) *ErrorService {
	statusCode := GetStatusCode(NO_CONTENT)
	stackTrace := captureStackTrace(message)
	return &ErrorService{
		StatusText:   http.StatusText(statusCode),
		StatusCode:   statusCode,
		ErrorCode:    string(code),
		ErrorMessage: message,
		StackTrace:   stackTrace,
	}
}

// ExtractErrorDetails extracts the details from ErrorService
func ExtractErrorDetails(err error) (res *ErrorService) {
	if err, ok := err.(*ErrorService); ok {
		return err
	}
	return res
}
