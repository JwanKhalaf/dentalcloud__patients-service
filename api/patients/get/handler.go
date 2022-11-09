package get

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jwankhalaf/dentalcloud__patients-service/api/patients"
	"go.uber.org/zap"
)

const contentTypeHeader string = "content-type"
const jsonContentType string = "application/json"

func GetPatientHandler(logger *zap.Logger, repository patients.PatientRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("running the get patient handler...")

		patientID := strings.TrimPrefix(r.URL.Path, "/patients/")

		logger.Info("requested patient", zap.String("patientID", patientID))

		logger = logger.With(zap.String("patientID", patientID))

		w.Header().Set(contentTypeHeader, jsonContentType)

		patient, err := repository.GetPatient(logger, r.Context(), patientID)
		if err != nil {
			logger.Error("failed to get patient", zap.Error(err))
			http.Error(w, "requested patient could not be found", http.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(patient)
		if err != nil {
			logger.Error("error in json marshal", zap.Error(err))
		}

		w.WriteHeader(http.StatusOK)
	})
}
