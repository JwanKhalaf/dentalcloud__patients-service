package create

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePatient(t *testing.T) {
	t.Run("returns 415 (unsupported media type) when the request does not have content-type set as application/json", func(t *testing.T) {
		// create a request to pass to our handler
		req, _ := http.NewRequest("POST", "/patients", nil)

		// set the content type
		req.Header.Set("content-type", "text/csv")

		// create a response recorder
		res := httptest.NewRecorder()

		// get the handler
		handler := CreatePatientHandler()

		// our handler satisfies http.handler, so we can call its serve http method
		// directly and pass in our request and response recorder
		handler.ServeHTTP(res, req)

		// assert status code is what we expect
		assertStatusCode(t, res.Code, http.StatusUnsupportedMediaType)
	})

	t.Run("create returns 400 (bad request) when first name is not set in the body", func(t *testing.T) {

		jsonValue, _ := json.Marshal(CreatePatientRequest{})

		// create a request to pass to our handler
		req, _ := http.NewRequest("POST", "/patients", bytes.NewBuffer(jsonValue))

		// set the content type
		req.Header.Set("content-type", "application/json")

		// create a response recorder
		res := httptest.NewRecorder()

		// get the handler
		handler := CreatePatientHandler()

		// our handler satisfies http.handler, so we can call its serve http method
		// directly and pass in our request and response recorder
		handler.ServeHTTP(res, req)

		// assert status code is what we expect
		assertStatusCode(t, res.Code, http.StatusBadRequest)
	})
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("handler returned wrong status code: got %v want %v", got, want)
	}
}
