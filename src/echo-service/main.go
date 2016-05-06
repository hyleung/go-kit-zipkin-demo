package main

import (
	echoservice "echo-service/echoservice"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting server...")
	ctx := context.Background()
	svc := echoservice.EchoServiceImpl{}
	handler := httptransport.NewServer(
		ctx,
		echoservice.MakeEchoEndpoint(svc),
		echoservice.DecodeEchoRequest,
		echoservice.EncodeEchoResponse,
	)
	http.Handle("/echo", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
