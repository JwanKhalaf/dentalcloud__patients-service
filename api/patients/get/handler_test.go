package get

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jwankhalaf/dentalcloud__patients-service/api/patients"
)

type StubPatientStore struct {
	getPatient func(ctx context.Context, patientID string) (patients.Patient, error)
}

func (s *StubPatientStore) GetPatient(ctx context.Context, patientID string) (patients.Patient, error) {
	return s.getPatient(ctx, patientID)
}

func TestGetPatient(t *testing.T) {
	t.Run("return 404 when requested patient does not exist", func(t *testing.T) {
		requestedPatientID := "test_patient_id"

		// create the stub patient store
		patientStore := StubPatientStore{
			getPatient: func(_ context.Context, patientID string) (patients.Patient, error) {
				if patientID != requestedPatientID {
					t.Errorf("%q was passed to GetPatient() but the expected value was %q", patientID, requestedPatientID)
				}

				return patients.Patient{}, errors.New("requested patient could not be found")
			},
		}

		// create a request to pass to our handler
		req, _ := http.NewRequest("GET", fmt.Sprintf("/patients/%v", requestedPatientID), nil)

		// create a response recorder
		res := httptest.NewRecorder()

		handler := GetPatientHandler(&patientStore)

		handler.ServeHTTP(res, req)

		// assert status code is what we expect
		assertStatusCode(t, res.Code, http.StatusNotFound)
	})

	t.Run("return 200 along with the patient details when requested patient does exist", func(t *testing.T) {
		requestedPatientID := "test_patient_id"
		expectedPatient := patients.Patient{PatientID: "test_patient_id", FirstName: "Jane", LastName: "Doe"}
		// create the stub patient store
		patientStore := StubPatientStore{
			getPatient: func(_ context.Context, patientID string) (patients.Patient, error) {
				if patientID != requestedPatientID {
					t.Errorf("%q was passed to GetPatient() but the expected value was %q", patientID, requestedPatientID)
				}

				return expectedPatient, nil
			},
		}

		// create a request to pass to our handler
		req, _ := http.NewRequest("GET", fmt.Sprintf("/patients/%v", requestedPatientID), nil)

		// create a response recorder
		res := httptest.NewRecorder()

		// get handler
		handler := GetPatientHandler(&patientStore)

		// our handler satisfies http.handler, so we can call its serve http method
		// directly and pass in our request and response recorder
		handler.ServeHTTP(res, req)

		// decode the json response into []tasks.Task
		got := getPatientFromResponse(t, res.Body)

		// check the response body is what we expect
		assertPatient(t, got, expectedPatient)

		// assert status code is what we expect
		assertStatusCode(t, res.Code, http.StatusOK)
	})
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("handler returned wrong status code: got %v want %v", got, want)
	}
}

func getPatientFromResponse(t testing.TB, body io.Reader) (patient patients.Patient) {
	t.Helper()

	err := json.NewDecoder(body).Decode(&patient)

	if err != nil {
		t.Fatalf("unable to process response from server %q into a Patient, '%v'", body, err)
	}

	return
}

func assertPatient(t testing.TB, got, want patients.Patient) {
	t.Helper()

	if diff := cmp.Diff(got, want); diff != "" {
		t.Error("handler returned unexpected body", diff)
	}
}
