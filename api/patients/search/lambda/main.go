package main

import (
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/jwankhalaf/dentalcloud__patients-service/api/patients"
	"github.com/jwankhalaf/dentalcloud__patients-service/api/patients/search"
	"go.uber.org/zap"
)

func main() {
	// initialise a new zap logger
	logger, _ := zap.NewProduction()

	logger.Info("running the search patients lamdba...")

	mux := http.NewServeMux()

	mux.Handle("/", search.SearchPatientsHandler(logger, patients.NewPatientStore(logger)))
	algnhsa.ListenAndServe(mux, nil)
}
