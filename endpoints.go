package main

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/mchudgins/jeeves/pkg/service"
	"github.com/mchudgins/jeeves/pkg/transport"
	"golang.org/x/net/context"
)

/*
  EndpointItem defines the API URL
*/
type EndpointItem struct {
	url     string
	e       endpoint.Endpoint
	decoder httptransport.DecodeRequestFunc
	encoder httptransport.EncodeResponseFunc
	options []httptransport.ServerOption
}

var (
	build = service.Build{}
	apis  = []EndpointItem{
		{"/test", transport.MakeLaunchBuildEndpoint(build), transport.DecodeBuildRequest, transport.EncodeResponse, nil},
	}
)

func createServer(e EndpointItem) *httptransport.Server {

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

	return handler
}

func CreateHttpEndpoints() {
	for _, e := range apis {
		http.Handle(e.url, createServer(e))
	}
}
