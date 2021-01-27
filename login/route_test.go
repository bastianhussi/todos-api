package login_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/bastianhussi/todos-api"
	"github.com/bastianhussi/todos-api/login"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name           string
		shouldFail     bool
		in             *http.Request
		out            *httptest.ResponseRecorder
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "pass",
			shouldFail:     false,
			in:             httptest.NewRequest(http.MethodPost, "/login", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusCreated,
			expectedBody:   "Hello, World!",
		},
		{
			name:           "fail",
			shouldFail:     true,
			in:             httptest.NewRequest(http.MethodGet, "/login", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello, World!",
		},
		{

			name:           "pass",
			shouldFail:     false,
			in:             httptest.NewRequest(http.MethodGet, "/login", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, err := api.NewConfig()
			must(err)
			res, err := api.NewResources(c)
			must(err)

			h := login.NewHandler(res)
			h.Login(test.out, test.in)

			code, body := test.out.Code, test.out.Body.String()

			if code != test.expectedStatus {
				if !test.shouldFail {
					t.Logf("expected: %d\ngot: %d\n", test.expectedStatus, code)
					t.Fail()
				}
			}

			if body != test.expectedBody {
				if !test.shouldFail {
					t.Logf("expected: %v\ngot: %v\n", test.expectedBody, body)
					t.Fail()
				}
			}
		})
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
