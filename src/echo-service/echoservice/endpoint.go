package echoservice

import (
	"github.com/go-kit/kit/endpoint"

	"golang.org/x/net/context"
)

func makeEchoEndpoint(svc EchoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(echoRequest)
		result := svc.Echo(req.Msg)
		return echoResponse{result}, nil
	}
}
