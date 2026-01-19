package jsonrpc

import (
	"fmt"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}

var (
	ErrParseError     = &Error{Code: -32700, Message: "Parse error"}
	ErrInvalidRequest = &Error{Code: -32600, Message: "Invalid request"}
	ErrMethodNotFound = &Error{Code: -32601, Message: "Method not found"}
	ErrInvalidParams  = &Error{Code: -32602, Message: "Invalid params"}
	ErrInternalError  = &Error{Code: -32603, Message: "Internal error"}
)
