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

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	var (
		port         = fs.Int("port", 8080, "Http port")
		scribeHost   = fs.String("scribeHost", "", "Scribe host <hostname:port>")
		samplingRate = fs.Float64("samplingRate", 1.0, "Sampling rate")
	)
	flag.Usage = fs.Usage
	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	fmt.Println("Starting server on port:", *port)
	//zipkin tracing
	serviceName := "echo-service"
	methodName := "echo"
	sampler := zipkin.SampleRate(*samplingRate, 0)
	spanFunc := zipkin.MakeNewSpanFunc(fmt.Sprintf("127.0.0.1:%d", *port), serviceName, methodName)
	timeout := time.Second
	batchInterval := time.Millisecond
	fmt.Println("Using Scribe collector at", *scribeHost, "sampling at", *samplingRate*100, "%")
	collector, err := zipkin.NewScribeCollector(*scribeHost, timeout,
		zipkin.ScribeBatchSize(0),
		zipkin.ScribeBatchInterval(batchInterval),
		zipkin.ScribeSampleRate(sampler))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	svc := echoservice.EchoServiceImpl{collector}
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
