package get

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/jwankhalaf/dentalcloud__patients-service/api/patients"
)

const contentTypeHeader string = "content-type"
const jsonContentType string = "application/json"

func GetPatientHandler(repository patients.PatientRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("running the get patient handler!")

		patientID := strings.TrimPrefix(r.URL.Path, "/patients/")

		log.Printf("requested patient id is %q", patientID)

		w.Header().Set(contentTypeHeader, jsonContentType)

		patient, err := repository.GetPatient(r.Context(), patientID)
		if err != nil {
			log.Printf("failed to get patient: %v", err)
			http.Error(w, "requested patient could not be found", http.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(patient)
		if err != nil {
			log.Printf("error in json marshal: %v", err)
		}

		w.WriteHeader(http.StatusOK)
	})
}
