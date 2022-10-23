package main

import (
	"log"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/jwankhalaf/dentalcloud__patients-service/api/patients"
	"github.com/jwankhalaf/dentalcloud__patients-service/api/patients/get"
)

func main() {
	log.Println("running the get patients lambda...")

	mux := http.NewServeMux()

	mux.Handle("/", get.GetPatientHandler(patients.NewPatientStore()))
	algnhsa.ListenAndServe(mux, nil)
}
