package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
)

// DecodeSumRequest decodes the request from the provided HTTP request, simply
// by JSON decoding from the request body. It's designed to be used in
// transport/http.Server.
func DecodeSumRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request SumRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	return &request, err
}

// EncodeSumResponse encodes the response to the provided HTTP response
// writer, simply by JSON encoding to the writer. It's designed to be used in
// transport/http.Server.
func EncodeSumResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// DecodeConcatRequest decodes the request from the provided HTTP request,
// simply by JSON decoding from the request body. It's designed to be used in
// transport/http.Server.
func DecodeConcatRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request ConcatRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	return &request, err
}

// EncodeConcatResponse encodes the response to the provided HTTP response
// writer, simply by JSON encoding to the writer. It's designed to be used in
// transport/http.Server.
func EncodeConcatResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// EncodeSumRequest encodes the request to the provided HTTP request, simply
// by JSON encoding to the request body. It's designed to be used in
// transport/http.Client.
func EncodeSumRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// DecodeSumResponse decodes the response from the provided HTTP response,
// simply by JSON decoding from the response body. It's designed to be used in
// transport/http.Client.
func DecodeSumResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response SumResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

// EncodeConcatRequest encodes the request to the provided HTTP request,
// simply by JSON encoding to the request body. It's designed to be used in
// transport/http.Client.
func EncodeConcatRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// DecodeConcatResponse decodes the response from the provided HTTP response,
// simply by JSON decoding from the response body. It's designed to be used in
// transport/http.Client.
func DecodeConcatResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response ConcatResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}
