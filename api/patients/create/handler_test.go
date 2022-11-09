package create

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePatient(t *testing.T) {
	t.Run("create returns 400 when missing required attributes", func(t *testing.T) {
		// create a request to pass to our handler
		req, _ := http.NewRequest("POST", "/patients", nil)

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
