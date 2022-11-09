package main

import (
	"log"
	"net/http"

	"github.com/akrylysov/algnhsa"
	"github.com/jwankhalaf/dentalcloud__patients-service/api/patients"
	"github.com/jwankhalaf/dentalcloud__patients-service/api/patients/search"
)

func main() {
	log.Println("running the search patients lambda...")

	mux := http.NewServeMux()

	mux.Handle("/", search.SearchPatientsHandler(patients.NewPatientStore()))
	algnhsa.ListenAndServe(mux, nil)
}
