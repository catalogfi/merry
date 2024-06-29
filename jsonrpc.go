package merry

import (
	"encoding/json"
)

// Error codes
const (
	ErrorCodeParseError        = -32700
	ErrorMessageParseError     = "Parse error"
	ErrorCodeInvalidRequest    = -32600
	ErrorMessageInvalidRequest = "Invalid Request"
	ErrorCodeMethodNotFound    = -32601
	ErrorMessageMethodNotFound = "Method not found"
	ErrorCodeInvalidParams     = -32602
	ErrorMessageInvalidParams  = "Invalid params"
	ErrorCodeInternalError     = -32603
	ErrorMessageInternalError  = "Internal error"
)

// Request defines a JSON-RPC 2.0 request object.
type Request struct {
	Version string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response defines a JSON-RPC 2.0 response object.
type Response struct {
	Version string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
}

// Error defines a JSON-RPC 2.0 error object.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

// NewResponse returns a new JSON-RPC 2.0 response object.
func NewResponse(id interface{}, result json.RawMessage, err *Error) Response {
	return Response{
		Version: "2.0",
		ID:      id,
		Result:  result,
		Error:   err,
	}
}

// NewError returns a new JSON-RPC 2.0 error object.
func NewError(code int, message string, data string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewParseError(err error) *Error {
	return &Error{
		Code:    ErrorCodeParseError,
		Message: ErrorMessageParseError,
		Data:    err.Error(),
	}
}

func NewInvalidRequest(err error) *Error {
	return &Error{
		Code:    ErrorCodeInvalidRequest,
		Message: ErrorMessageInvalidRequest,
		Data:    err.Error(),
	}
}

func NewInvalidParams(err error) *Error {
	return &Error{
		Code:    ErrorCodeInvalidParams,
		Message: ErrorMessageInvalidParams,
		Data:    err.Error(),
	}
}

func NewMethodNotFound() *Error {
	return &Error{
		Code:    ErrorCodeMethodNotFound,
		Message: ErrorMessageMethodNotFound,
		Data:    "",
	}
}
