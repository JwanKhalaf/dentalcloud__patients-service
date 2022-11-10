package create

import (
	"encoding/json"
	"mime"
	"net/http"
)

type CreatePatientRequest struct {
	FirstName string `json:"first_name"`
}

func CreatePatientHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// enforce a json content-type
		contentType := r.Header.Get("content-type")

		mediatype, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if mediatype != "application/json" {
			http.Error(w, "api expects application/json content-type", http.StatusUnsupportedMediaType)
			return
		}

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var createPatientRequest CreatePatientRequest
		if err := dec.Decode(&createPatientRequest); err != nil {
			http.Error(w, "request body is invalid", http.StatusBadRequest)
			return
		}

		if createPatientRequest.FirstName == "" {
			http.Error(w, "request body is invalid", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}
