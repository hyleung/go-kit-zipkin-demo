package echoservice

import (
	zipkin "github.com/go-kit/kit/tracing/zipkin"
	"golang.org/x/net/context"
	"strings"
)

//EchoService takes an input string returns it, uppercased.
type EchoService interface {
	Echo(context.Context, string) string
}

type EchoServiceImpl struct {
	Collector zipkin.Collector
}

func (svc EchoServiceImpl) Echo(ctx context.Context, s string) string {
	span, collect := zipkin.NewChildSpan(
		ctx,
		svc.Collector,
		"echoing",
		zipkin.ServerAddr(
			"127.0.0.1",
			"echoserver",
		),
	)
	span.AnnotateBinary("Echo param", s)
	collect()
	return strings.ToUpper(s)
}
