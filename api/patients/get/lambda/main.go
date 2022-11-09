package main

import (
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/jwankhalaf/dentalcloud__patients-service/api/patients"
	"github.com/jwankhalaf/dentalcloud__patients-service/api/patients/get"
	"go.uber.org/zap"
)

func main() {
	// initialise a new zap logger
	logger, _ := zap.NewProduction()

	logger.Info("running the get patient lamdba...")

	mux := http.NewServeMux()

	mux.Handle("/", get.GetPatientHandler(logger, patients.NewPatientStore(logger)))
	algnhsa.ListenAndServe(mux, nil)
}
