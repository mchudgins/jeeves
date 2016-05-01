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
	Name string
}

type BuildService interface {
	LaunchBuild(ctx context.Context, build BuildInfo) (string, error)
}

func (svc Build) LaunchBuild(ctx context.Context, build BuildInfo) (string, error) {
	txID, _ := CorrelationIDFromContext(ctx)
	fmt.Printf("LaunchBuild (txID: %s)\n", txID)
	return build.Name, nil
}

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

// userIPkey is the context key for the user IP address.  Its value of zero is
// arbitrary.  If this package defined other context keys, they would have
// different integer values.
const correlationIDKey key = 0

func NewContext(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationIDKey, correlationID)
}

func CorrelationIDFromContext(ctx context.Context) (string, bool) {
	// ctx.Value returns nil if ctx has no value for the key;
	// the string type assertion returns ok=false for nil.
	id, ok := ctx.Value(correlationIDKey).(string)
	return id, ok
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
