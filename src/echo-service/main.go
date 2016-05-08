package main

import (
	echoservice "echo-service/echoservice"
	"flag"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	zipkin "github.com/go-kit/kit/tracing/zipkin"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"time"
)

var port = flag.Int("port", 8080, "Http port")
var scribeHost = flag.String("scribeHost", "localhost:9410", "Scribe host <hostname:port>")

func main() {
	flag.Parse()
	fmt.Println("Starting server on port:", *port)
	//zipkin tracing
	serviceName := "echo-service"
	methodName := "echo"

	spanFunc := zipkin.MakeNewSpanFunc(fmt.Sprintf("127.0.0.1:%d", *port), serviceName, methodName)
	timeout := time.Second
	batchInterval := time.Millisecond
	fmt.Println("Using Scribe collector at ", *scribeHost)
	collector, err := zipkin.NewScribeCollector(*scribeHost, timeout, zipkin.ScribeBatchSize(0), zipkin.ScribeBatchInterval(batchInterval))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	svc := echoservice.EchoServiceImpl{}
	var service endpoint.Endpoint
	service = echoservice.MakeEchoEndpoint(svc)
	//wrap the service in Zipkin tracing middleware
	service = zipkin.AnnotateServer(spanFunc, collector)(service)

	handler := httptransport.NewServer(
		ctx,
		service,
		echoservice.DecodeEchoRequest,
		echoservice.EncodeEchoResponse,
		httptransport.ServerBefore(
			zipkin.ToContext(spanFunc, kitlog.NewLogfmtLogger(os.Stdout)),
		),
	)
	http.Handle("/echo", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
