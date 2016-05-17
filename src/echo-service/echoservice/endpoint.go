package echoservice

import (
	"github.com/go-kit/kit/endpoint"

	"golang.org/x/net/context"
)

func MakeEchoEndpoint(svc EchoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(echoRequest)
		result := svc.Echo(ctx, req.Msg)
		return echoResponse{result}, nil
	}
}
