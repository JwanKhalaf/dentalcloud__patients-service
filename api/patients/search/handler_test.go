package search

import (
	"context"
	"encoding/json"
	"fmt"
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

func TestSearchPatient(t *testing.T) {
	// create the logger
	logger, _ := zap.NewProduction()

	t.Run("returns a bad request if endpoint is called without search query param", func(t *testing.T) {
		var p = make([]patients.PatientSearchResponseItem, 0)

		// create the stub patient store
		patientStore := StubPatientStore{
			searchPatients: func(_ *zap.Logger, _ context.Context, searchTerm string) ([]patients.PatientSearchResponseItem, error) {
				return p, nil
			},
		}

		// create a request to pass to our handler
		req, _ := http.NewRequest("GET", "/patients", nil)

		// create a response recorder
		res := httptest.NewRecorder()

		// get the handler
		handler := SearchPatientsHandler(logger, &patientStore)

		// our handler satisfies http.handler, so we can call its serve http method
		// directly and pass in our request and response recorder
		handler.ServeHTTP(res, req)

		// assert status code is what we expect
		assertStatusCode(t, res.Code, http.StatusBadRequest)
	})

	t.Run("returns empty list when search param has no value", func(t *testing.T) {
		searchParam := ""
		var p = make([]patients.PatientSearchResponseItem, 0)

		// create the stub patient store
		patientStore := StubPatientStore{
			searchPatients: func(_ *zap.Logger, _ context.Context, searchTerm string) ([]patients.PatientSearchResponseItem, error) {
				return p, nil
			},
		}

		// create a request to pass to our handler
		req, _ := http.NewRequest("GET", fmt.Sprintf("/patients?search=%v", searchParam), nil)

		// create a response recorder
		res := httptest.NewRecorder()

		// get the handler
		handler := SearchPatientsHandler(logger, &patientStore)

		// our handler satisfies http.handler, so we can call its serve http method
		// directly and pass in our request and response recorder
		handler.ServeHTTP(res, req)

		// decode the json response into []patients.PatientSearchResponseItem
		got := getPatientsFromResponse(t, res.Body)

		// assert status code is what we expect
		assertStatusCode(t, res.Code, http.StatusOK)

		// assert response is an epty json list
		assertSearchResponse(t, got, p)
	})

	t.Run("check that search param is getting passed to the patients store", func(t *testing.T) {
		searchParam := "jam"

		// create the stub patient store
		patientStore := StubPatientStore{
			searchPatients: func(_ *zap.Logger, _ context.Context, searchTerm string) ([]patients.PatientSearchResponseItem, error) {
				if searchTerm != searchParam {
					t.Errorf("%q was passed to SearchPatients() but the expected value was %q", searchTerm, searchParam)
				}

				return []patients.PatientSearchResponseItem{}, nil
			},
		}

		// create a request to pass to our handler
		req, _ := http.NewRequest("GET", fmt.Sprintf("/patients?search=%v", searchParam), nil)

		// create a response recorder
		res := httptest.NewRecorder()

		// get the handler
		handler := SearchPatientsHandler(logger, &patientStore)

		// our handler satisfies http.handler, so we can call its serve http method
		// directly and pass in our request and response recorder
		handler.ServeHTTP(res, req)

		// decode the json response into []patients.PatientSearchResponseItem
		getPatientsFromResponse(t, res.Body)
	})

	t.Run("check response is correct when an array of patients is returned", func(t *testing.T) {
		searchParam := "jam"

		expectedPatients := []patients.PatientSearchResponseItem{
			{PatientID: "test_patient_id_1", FirstName: "jamie", LastName: "oliver", DateOfBirth: "test_dob", Email: "j.oliver@gmail.com", MobilePhone: "07865154788", PostCode: "LS18 9BQ"},
			{PatientID: "test_patient_id_2", FirstName: "james", LastName: "watt", DateOfBirth: "test_dob", Email: "j.watt@gmail.com", MobilePhone: "07531247866", PostCode: "LS1 3LP"},
		}

		// create the stub patient store
		patientStore := StubPatientStore{
			searchPatients: func(_ *zap.Logger, _ context.Context, searchTerm string) ([]patients.PatientSearchResponseItem, error) {
				return expectedPatients, nil
			},
		}

		// create a request to pass to our handler
		req, _ := http.NewRequest("GET", fmt.Sprintf("/patients?search=%v", searchParam), nil)

		// create a response recorder
		res := httptest.NewRecorder()

		// get the handler
		handler := SearchPatientsHandler(logger, &patientStore)

		// our handler satisfies http.handler, so we can call its serve http method
		// directly and pass in our request and response recorder
		handler.ServeHTTP(res, req)

		// decode the json response into []patients.PatientSearchResponseItem
		got := getPatientsFromResponse(t, res.Body)

		// assert status code is what we expect
		assertStatusCode(t, res.Code, http.StatusOK)

		// assert response is an epty json list
		assertSearchResponse(t, got, expectedPatients)
	})
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("handler returned wrong status code: got %v want %v", got, want)
	}
}

func getPatientsFromResponse(t testing.TB, body io.Reader) (patients []patients.PatientSearchResponseItem) {
	t.Helper()

	err := json.NewDecoder(body).Decode(&patients)

	if err != nil {
		t.Fatalf("unable to process response from server %q into a Patient, '%v'", body, err)
	}

	return
}

func assertSearchResponse(t testing.TB, got, want []patients.PatientSearchResponseItem) {
	t.Helper()

	if diff := cmp.Diff(got, want); diff != "" {
		t.Error("handler returned unexpected body", diff)
	}
}
