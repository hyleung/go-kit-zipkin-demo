package main

import (
	echoservice "echo-service/echoservice"
	"flag"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
	"log"
	"net/http"
)

var port = flag.Int("port", 8080, "Http port")

func main() {
	flag.Parse()
	fmt.Println("Starting server on port:", *port)
	ctx := context.Background()
	svc := echoservice.EchoServiceImpl{}
	handler := httptransport.NewServer(
		ctx,
		echoservice.MakeEchoEndpoint(svc),
		echoservice.DecodeEchoRequest,
		echoservice.EncodeEchoResponse,
	)
	http.Handle("/echo", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
