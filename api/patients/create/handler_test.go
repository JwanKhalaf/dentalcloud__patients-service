package create

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jwankhalaf/dentalcloud__patients-service/api/patients"
	"go.uber.org/zap"
)

type StubPatientStore struct {
	createPatient  func(logger *zap.Logger, ctx context.Context, patient patients.CreatePatientRequest) (patients.CreatePatientResponse, error)
	getPatient     func(logger *zap.Logger, ctx context.Context, patientID string) (patients.Patient, error)
	searchPatients func(logger *zap.Logger, ctx context.Context, searchTerm string) ([]patients.PatientSearchResponseItem, error)
}

func (s *StubPatientStore) CreatePatient(logger *zap.Logger, ctx context.Context, patient patients.CreatePatientRequest) (patients.CreatePatientResponse, error) {
	return s.createPatient(logger, ctx, patient)
}

func (s *StubPatientStore) GetPatient(logger *zap.Logger, ctx context.Context, patientID string) (patients.Patient, error) {
	return s.getPatient(logger, ctx, patientID)
}

func (s *StubPatientStore) SearchPatients(logger *zap.Logger, ctx context.Context, searchTerm string) ([]patients.PatientSearchResponseItem, error) {
	return s.searchPatients(logger, ctx, searchTerm)
}

func TestCreatePatient(t *testing.T) {
	contentType := "content-type"
	applicationJson := "application/json"

	// create the logger
	logger, _ := zap.NewProduction()

	t.Run("returns 415 (unsupported media type) when the request does not have content-type set as application/json", func(t *testing.T) {
		patientToBeCreated := patients.CreatePatientRequest{FirstName: "James"}

		// create the stub patient store
		patientsStore := StubPatientStore{
			createPatient: func(_ *zap.Logger, _ context.Context, patient patients.CreatePatientRequest) (patients.CreatePatientResponse, error) {
				if patient.FirstName != patientToBeCreated.FirstName {
					t.Errorf("got: CreatePatient(%s) expected CreatePatient(%s)", patient.FirstName, patientToBeCreated.FirstName)
				}
				return patients.CreatePatientResponse{}, nil
			},
		}

		// create a request to pass to our handler
		req, _ := http.NewRequest("POST", "/patients", nil)

		// set the content type
		req.Header.Set(contentType, "text/csv")

		// create a response recorder
		res := httptest.NewRecorder()

		// get the handler
		handler := CreatePatientHandler(logger, &patientsStore)

		// our handler satisfies http.handler, so we can call its serve http method
		// directly and pass in our request and response recorder
		handler.ServeHTTP(res, req)

		// assert status code is what we expect
		assertStatusCode(t, res.Code, http.StatusUnsupportedMediaType)
	})

	t.Run("create returns 400 (bad request) when first name is not set in the body", func(t *testing.T) {
		// the patient to be created
		patientToBeCreated := patients.CreatePatientRequest{}

		// create the stub patient store
		patientsStore := StubPatientStore{
			createPatient: func(_ *zap.Logger, _ context.Context, patient patients.CreatePatientRequest) (patients.CreatePatientResponse, error) {
				if patient.FirstName != patientToBeCreated.FirstName {
					t.Errorf("got: CreatePatient(%s) expected CreatePatient(%s)", patient.FirstName, patientToBeCreated.FirstName)
				}
				return patients.CreatePatientResponse{}, nil
			},
		}

		jsonValue, _ := json.Marshal(patientToBeCreated)

		// create a request to pass to our handler
		req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer(jsonValue))

		// set the content type
		req.Header.Set(contentType, applicationJson)

		// create a response recorder
		res := httptest.NewRecorder()

		// get the handler
		handler := CreatePatientHandler(logger, &patientsStore)

		// our handler satisfies http.handler, so we can call its serve http method
		// directly and pass in our request and response recorder
		handler.ServeHTTP(res, req)

		// assert status code is what we expect
		assertStatusCode(t, res.Code, http.StatusBadRequest)
	})

	t.Run("create returns 500 (internal server error) when call to dynamodb fails", func(t *testing.T) {
		// the patient to be created
		patientToBeCreated := patients.CreatePatientRequest{FirstName: "Jason"}

		// create the stub patient store
		patientsStore := StubPatientStore{
			createPatient: func(_ *zap.Logger, _ context.Context, patient patients.CreatePatientRequest) (patients.CreatePatientResponse, error) {
				if patient.FirstName != patientToBeCreated.FirstName {
					t.Errorf("got: CreatePatient(%s) expected CreatePatient(%s)", patient.FirstName, patientToBeCreated.FirstName)
				}
				return patients.CreatePatientResponse{}, errors.New("call to dynamodb failed")
			},
		}

		jsonValue, _ := json.Marshal(patientToBeCreated)

		// create a request to pass to our handler
		req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer(jsonValue))

		// set the content type
		req.Header.Set(contentType, applicationJson)

		// create a response recorder
		res := httptest.NewRecorder()

		// get the handler
		handler := CreatePatientHandler(logger, &patientsStore)

		// our handler satisfies http.handler, so we can call its serve http method
		// directly and pass in our request and response recorder
		handler.ServeHTTP(res, req)

		// assert status code is what we expect
		assertStatusCode(t, res.Code, http.StatusInternalServerError)
	})

	t.Run("create returns 201 (created) and new patient id when first name is set in the body", func(t *testing.T) {
		// the patient to be created
		expectedPatient := patients.CreatePatientRequest{FirstName: "James"}

		// the expected handler response
		expectedResponse := patients.CreatePatientResponse{PatientID: "test_id"}

		// create the stub patient store
		patientsStore := StubPatientStore{
			createPatient: func(_ *zap.Logger, _ context.Context, patient patients.CreatePatientRequest) (patients.CreatePatientResponse, error) {
				if patient.FirstName != expectedPatient.FirstName {
					t.Errorf("got: CreatePatient(%s) expected CreatePatient(%s)", patient, expectedPatient)
				}
				return expectedResponse, nil
			},
		}

		jsonValue, _ := json.Marshal(expectedPatient)

		// create a request to pass to our handler
		req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer(jsonValue))

		// set the content type
		req.Header.Set(contentType, applicationJson)

		// create a response recorder
		res := httptest.NewRecorder()

		// get the handler
		handler := CreatePatientHandler(logger, &patientsStore)

		// our handler satisfies http.handler, so we can call its serve http method
		// directly and pass in our request and response recorder
		handler.ServeHTTP(res, req)

		// decode the json response into patients.CreatePatientResponse
		got := getResponse(t, res.Body)

		// assert status code is what we expect
		assertStatusCode(t, res.Code, http.StatusCreated)

		// assert response body
		assertCreatePatientResponse(t, got, expectedResponse)
	})
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("handler returned wrong status code: got %v want %v", got, want)
	}
}

func getResponse(t testing.TB, body io.Reader) (createPatientResponse patients.CreatePatientResponse) {
	t.Helper()

	err := json.NewDecoder(body).Decode(&createPatientResponse)

	if err != nil {
		t.Fatalf("unable to process response from server %q into a CreatePatientResponse, '%v'", body, err)
	}

	return
}

func assertCreatePatientResponse(t testing.TB, got, want patients.CreatePatientResponse) {
	t.Helper()

	if diff := cmp.Diff(got, want); diff != "" {
		t.Error("handler returned unexpected body", diff)
	}
}
