package main

import (
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/jwankhalaf/dentalcloud__patients-service/api/patients"
	"github.com/jwankhalaf/dentalcloud__patients-service/api/patients/create"
	"go.uber.org/zap"
)

func main() {
	// initialise a new zap logger
	logger, _ := zap.NewProduction()

	logger.Info("running the create patient lamdba...")

	mux := http.NewServeMux()

	mux.Handle("/", create.CreatePatientHandler(logger, patients.NewPatientStore(logger)))
	algnhsa.ListenAndServe(mux, nil)
}
