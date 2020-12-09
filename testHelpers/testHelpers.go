package testHelpers

import (
	"io/ioutil"
	"encoding/json"
	models "fantasy_league/Models"
	"fmt"
	"io"
	"os"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// can go to mocks?

func NewGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func NewPostScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func NewGetLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/league"), nil)
	return req
}

// could be abstracted further as assert-helpers maybe

func AssertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("Got %q want %q", got, want)
	}
}

func AssertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("Got %d want %d", got, want)
	}
}

func AssertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func AssertLeague(t *testing.T, got, want []models.Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func ParseLeagueFromResponse(t *testing.T, body io.Reader) (league []models.Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}
	return
}



func CreateTempFile(t *testing.T, initialData string)(*os.File, func()) {
	t.Helper()
	tmpFile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("Could not create temp file!! %v", err)
	}
	
	tmpFile.Write([]byte(initialData))

	removeFile := func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())

	}
	return tmpFile, removeFile
}

func AssertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("didn't expect an error but got one, %v", err)
    }
}