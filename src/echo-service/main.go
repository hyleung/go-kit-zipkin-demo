package main

import (
	echoservice "echo-service/echoservice"
	"flag"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	log "github.com/go-kit/kit/log"
	kitot "github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go-opentracing"
	"golang.org/x/net/context"
	"io"
	stdlog "log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	Trace *stdlog.Logger
	Info  *stdlog.Logger
	Error *stdlog.Logger
)

func Init(traceHandle io.Writer, infoHandle io.Writer, errorHandle io.Writer) {
	Trace = stdlog.New(traceHandle, "TRACE :", stdlog.Ldate|stdlog.Ltime|stdlog.Lshortfile)
	Info = stdlog.New(traceHandle, "INFO:", stdlog.Ldate|stdlog.Ltime|stdlog.Lshortfile)
	Error = stdlog.New(traceHandle, "ERROR:", stdlog.Ldate|stdlog.Ltime|stdlog.Lshortfile)
}
func main() {
	Init(os.Stdout, os.Stdout, os.Stderr)
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
	// package log
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC).With("caller", log.DefaultCaller)
		stdlog.SetFlags(0)                             // flags are handled by Go kit's logger
		stdlog.SetOutput(log.NewStdlibAdapter(logger)) // redirect anything using stdlib log to us
	}
	Info.Println("Starting server on port:", *port)
	//tracing
	var tracer opentracing.Tracer
	{
		switch {
		case *scribeHost != "":
			serviceName := "echo-service"
			timeout := time.Second
			Info.Println("Using Scribe collector at", *scribeHost, "sampling at", *samplingRate*100, "%")
			collector, err := zipkintracer.NewScribeCollector(*scribeHost, timeout, zipkintracer.ScribeBatchSize(1))

			if err != nil {
				stdlog.Fatal(err)
			}
			ipaddr, _ := getLocalIP()
			recorder := zipkintracer.NewRecorder(collector, true, fmt.Sprintf("%s:%d", ipaddr, *port), serviceName)
			sampler := zipkintracer.NewCountingSampler(*samplingRate)
			tracer, err = zipkintracer.NewTracer(recorder, zipkintracer.WithSampler(sampler), zipkintracer.WithLogger(zipkintracer.LogWrapper(Error)))
			if err != nil {
				stdlog.Fatal(err)
			}
		default:
			Info.Println("Defaulting to no-op tracer...")
			tracer = opentracing.GlobalTracer()
		}
	}

	ctx := context.Background()
	svc := echoservice.EchoServiceImpl{}
	var service endpoint.Endpoint
	service = echoservice.MakeEchoEndpoint(svc)
	var (
		transportLogger = log.NewContext(logger).With("transport", "HTTP/JSON")
		tracingLogger   = log.NewContext(transportLogger).With("echo", "tracing")
	)
	//wrap the service in Zipkin tracing middleware
	service = kitot.TraceServer(tracer, "echo")(service)
	Info.Println("Server tracing initialized...")
	handler := httptransport.NewServer(
		ctx,
		service,
		echoservice.DecodeEchoRequest,
		echoservice.EncodeEchoResponse,
		httptransport.ServerErrorLogger(transportLogger),
		httptransport.ServerBefore(kitot.FromHTTPRequest(tracer, "echo", tracingLogger)),
	)
	http.Handle("/echo", handler)
	stdlog.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func getLocalIP() (ipaddr string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("Unable to get local IP")

}
