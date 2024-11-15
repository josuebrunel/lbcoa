package apiresponse

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ApiResponse[T any] struct {
	Error string `json:"error"`
	Data  T      `json:"data"`
}

func New[T any](w http.ResponseWriter, status int, data T, err error) {
	if err == nil {
		err = errors.New("")
	}
	b, _ := json.Marshal(ApiResponse[T]{
		Error: err.Error(),
		Data:  data,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
}
