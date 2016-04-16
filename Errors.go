package main

import "fmt"

type ErrorCode struct {
	Code         int
	ErrorDetails error
}

func (e ErrorCode) Error() string {
	return fmt.Sprintf("code: %d; %v", e.Code, e.ErrorDetails)
}

func NewErrorCode(code int, err error) ErrorCode {
	return ErrorCode{Code: code, ErrorDetails: err}
}
