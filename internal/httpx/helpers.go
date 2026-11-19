package httpx

import (
	"encoding/json"
	"net/http"
)

type Envelope map[string]any

func WriteJSON(w http.ResponseWriter, data any, status int, header http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, values := range header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)

	return err
}
