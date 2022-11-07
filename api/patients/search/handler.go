package search

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jwankhalaf/dentalcloud__patients-service/api/patients"
)

const contentTypeHeader string = "content-type"
const jsonContentType string = "application/json"

func SearchPatientsHandler(repository patients.PatientRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("running the search patients handler...")

		w.Header().Set(contentTypeHeader, jsonContentType)

		v, exist := r.URL.Query()["search"]
		if !exist {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		searchResults, err := repository.SearchPatients(r.Context(), v[0])
		if err != nil {
			log.Printf("failed to search patients: %v", err)
			http.Error(w, "requested patient could not be found", http.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(searchResults)
		if err != nil {
			log.Printf("error in json marshal: %v", err)
		}

		w.WriteHeader(http.StatusOK)
	})
}
