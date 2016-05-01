package transport

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mchudgins/jeeves/pkg/service"
)

// CoreRequest contains the two fields every request should have:
// a correlation ID and a user ID.
type CoreRequest struct {
	txID   string
	userID string
}

func (c *CoreRequest) populateCoreRequest(r *http.Request) {
	if txid := r.Header.Get("X-Correlation-ID"); len(txid) > 0 {
		c.txID = txid
	}

	if remoteUser := r.Header.Get("X-Remote-User"); len(remoteUser) > 0 {
		c.userID = remoteUser
	}
}

type LaunchBuildRequest struct {
	CoreRequest
	Name string `json:name`
}

type LaunchBuildResponse struct {
	Name  string `json:name`
	Error string `json:"error,omitempty"`
}

func MakeLaunchBuildEndpoint(svc service.BuildService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LaunchBuildRequest)

		fmt.Printf("request: %+v\n", req)
		build := service.BuildInfo{Name: req.Name}

		if len(req.txID) > 0 {
			ctx = service.NewContext(ctx, req.txID)
		}

		n, err := svc.LaunchBuild(ctx, build)
		if err != nil {
			return LaunchBuildResponse{n, err.Error()}, nil
		}

		return LaunchBuildResponse{n, ""}, nil
	}
}

func DecodeEmptyRequest(r *http.Request) (interface{}, error) {
	return r, nil
}

func DecodeBuildRequest(r *http.Request) (interface{}, error) {
	req := LaunchBuildRequest{}

	req.populateCoreRequest(r)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func EncodeResponse(w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
