package service

import (
	"fmt"

	"golang.org/x/net/context"
)
import "time"

type Build struct {
	Build       string
	LastUpdated time.Time
	Namespace   string
	Number      int
}

type BuildInfo struct {
	name string
}

type BuildService interface {
	LaunchBuild(ctx context.Context, build BuildInfo) (string, error)
}

func (svc Build) LaunchBuild(ctx context.Context, build BuildInfo) (string, error) {
	return "", nil
}

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
