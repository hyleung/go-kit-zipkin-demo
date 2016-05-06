package echoservice

import (
	"encoding/json"

	"golang.org/x/net/context"

	"net/http"
)

type echoRequest struct {
	Msg string `json:"msg"`
}

type echoResponse struct {
	Echo string `json:"echo"`
}

func DecodeEchoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request echoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeEchoResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
