package transports

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mchudgins/jeeves/pkg/service"
)

type LaunchBuildResponse struct {
	Name  string `json:name`
	Error string `json:"error,omitempty"`
}

func makeLaunchBuildEndpoint(svc service.BuildService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.BuildInfo)

		n, err := svc.LaunchBuild(ctx, req)
		if err != nil {
			return LaunchBuildResponse{n, err.Error()}, nil
		}

		return LaunchBuildResponse{n, ""}, nil
	}
}

func DecodeEmptyRequest(r *http.Request) (interface{}, error) {
	return r, nil
}

func EncodeResponse(w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
