package transports

import (
	"encoding/json"
	"net/http"
)

func DecodeEmptyRequest(r *http.Request) (interface{}, error) {
	return r, nil
}

func EncodeResponse(w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
