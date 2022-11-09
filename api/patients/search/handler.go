package search

import (
	"encoding/json"
	"net/http"

	"github.com/jwankhalaf/dentalcloud__patients-service/api/patients"
	"go.uber.org/zap"
)

const contentTypeHeader string = "content-type"
const jsonContentType string = "application/json"

func SearchPatientsHandler(logger *zap.Logger, repository patients.PatientRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("running the search patients handler...")

		w.Header().Set(contentTypeHeader, jsonContentType)

		v, exist := r.URL.Query()["search"]
		if !exist {
			logger.Error("no search term set as part of the query string params")
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		searchTerm := v[0]

		logger = logger.With(zap.String("searchTerm", searchTerm))

		searchResults, err := repository.SearchPatients(logger, r.Context(), searchTerm)
		if err != nil {
			logger.Error("failed to search patients", zap.Error(err))
			http.Error(w, "requested patient could not be found", http.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(searchResults)
		if err != nil {
			logger.Error("error in json marshal", zap.Error(err))
		}

		w.WriteHeader(http.StatusOK)
	})
}
