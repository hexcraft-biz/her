package her

import (
	"encoding/json"
	"errors"
	"net/http"
)

type errInterface struct {
	text *string
}

func (e errInterface) Error() string {
	if e.text != nil {
		return *e.text
	}
	return ""
}

// ================================================================
//
// ================================================================
type Payload struct {
	Message string `json:"message"`
	Result  any    `json:"result,omitempty"`
}

func NewPayload(result any) *Payload {
	return &Payload{
		Result: result,
	}
}

type Error struct {
	StatusCode int
	*Payload
	*errInterface
}

func (e Error) Is(target error) bool {
	return e == target
}

func (e Error) HttpR() (int, *Payload) {
	if e.StatusCode == http.StatusNoContent {
		return e.StatusCode, nil
	} else {
		return e.StatusCode, e.Payload
	}
}

// ================================================================
func New(code int, result any) *Error {
	e := &Error{
		StatusCode: code,
		Payload: &Payload{
			Message: http.StatusText(code),
			Result:  result,
		},
	}

	if code >= 400 {
		e.errInterface = &errInterface{text: &e.Message}
	}

	return e
}

// Return an *Error with err passing in. return nil if err is nil.
func NewError(code int, err error, result any) *Error {
	var her *Error

	if err != nil {
		her = NewErrorWithMessage(code, err.Error(), result)
	}

	return her
}

// ================================================================
func NewErrorWithMessage(code int, msg string, result any) *Error {
	her := &Error{
		StatusCode: code,
		Payload: &Payload{
			Message: msg,
			Result:  result,
		},
	}

	her.errInterface = &errInterface{text: &her.Payload.Message}
	return her
}

var (
	ErrBadRequest          = NewErrorWithMessage(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), nil)
	ErrUnauthorized        = NewErrorWithMessage(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
	ErrForbidden           = NewErrorWithMessage(http.StatusForbidden, http.StatusText(http.StatusForbidden), nil)
	ErrNotFound            = NewErrorWithMessage(http.StatusNotFound, http.StatusText(http.StatusNotFound), nil)
	ErrConflict            = NewErrorWithMessage(http.StatusConflict, http.StatusText(http.StatusConflict), nil)
	ErrGone                = NewErrorWithMessage(http.StatusGone, http.StatusText(http.StatusGone), nil)
	ErrUnprocessableEntity = NewErrorWithMessage(http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity), nil)
	ErrServiceUnavailable  = NewErrorWithMessage(http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable), nil)
	ErrInternalServerError = NewErrorWithMessage(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), nil)

	Errs = errors.Join(
		ErrBadRequest,
		ErrUnauthorized,
		ErrForbidden,
		ErrNotFound,
		ErrConflict,
		ErrGone,
		ErrUnprocessableEntity,
		ErrServiceUnavailable,
		ErrInternalServerError,
	)
)

// ================================================================
//
// ================================================================
func Assert(err error) *Error {
	if her, ok := err.(*Error); ok {
		return her
	} else {
		return nil
	}
}

func FetchHexcApiResult(resp *http.Response, payload *Payload) *Error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(payload); err != nil {
		return NewError(http.StatusInternalServerError, err, nil)
	}

	if resp.StatusCode >= 500 {
		return NewErrorWithMessage(http.StatusServiceUnavailable, payload.Message, nil)
	}

	return nil
}
