package request

import (
	"net/http"
)

func HandleBody[T any](r *http.Request) (*T, error) {
	payload, err := Decode[T](r.Body)
	if err != nil {
		return nil, err
	}
	err = IsValid(payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
