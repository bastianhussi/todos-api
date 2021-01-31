package profile_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var server *httptest.Server

func setup() {
	server := httptest.NewServer(&myHandler)
}

func teardown() {
	server.Close()
}

func TestSomething(t *testing.T) {
	setup()
	defer teardown()

	r, err := http.NewRequest(...)

	res, err := http.DefaultClient.Do(r)

	body, err := ioutil.ReadAll(res.body)
}

func TestAnother(t *testing.T) {
	
}
