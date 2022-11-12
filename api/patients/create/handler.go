package create

import (
	"encoding/json"
	"mime"
	"net/http"

	"github.com/jwankhalaf/dentalcloud__patients-service/api/patients"
	"go.uber.org/zap"
)

func CreatePatientHandler(logger *zap.Logger, repository patients.PatientRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// enforce a json content-type
		contentType := r.Header.Get("content-type")

		mediatype, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			logger.Error("error when parsing the mime type", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if mediatype != "application/json" {
			logger.Error("unsupported content-type", zap.Error(err))
			http.Error(w, "api expects application/json content-type", http.StatusUnsupportedMediaType)
			return
		}

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		var createPatientRequest patients.CreatePatientRequest
		if err := dec.Decode(&createPatientRequest); err != nil {
			logger.Error("the request body is invalid", zap.Error(err))
			http.Error(w, "request body is invalid", http.StatusBadRequest)
			return
		}

		// validation
		if createPatientRequest.FirstName == "" {
			http.Error(w, "request body is invalid", http.StatusBadRequest)
			return
		}

		response, err := repository.CreatePatient(logger, r.Context(), createPatientRequest)
		if err != nil {
			logger.Error("failed to create the patient", zap.Error(err))
			http.Error(w, "failed to create the patient", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			logger.Error("failed to encode the json for the create patient response", zap.Error(err))
		}
	})
}
