package main

import (
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	http "github.com/mchudgins/jeeves/pkg/transport"
	"golang.org/x/net/context"
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

func CreateServer(e EndpointItem) (*httptransport.Server, error) {

	ctx := context.Background()

	handler := httptransport.NewServer(
		ctx,
		//		endpoint.Chain(endpointInstrumentation(&counters, "createCertificate"),
		//			endpointLog("createCertificate"))(makeCreateCertificateEndpoint(svc)),
		e.e,
		e.decoder,
		e.encoder,
		e.options...,
	)

	return handler, nil
}
