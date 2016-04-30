package main

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	http "github.com/mchudgins/jeeves/pkg/transport"
)

/*
  EndpointItem defines the API URL
*/
type EndpointItem struct {
	e       endpoint.Endpoint
	decoder httptransport.DecodeRequestFunc
	encoder httptransport.EncodeResponseFunc
	options []httptransport.ServerOption
}

var (
	apis = []EndpointItem{
		{nil, http.DecodeEmptyRequest, http.EncodeResponse, nil},
	}
)
