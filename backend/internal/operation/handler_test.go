package operation

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func request(t *testing.T, handler http.Handler, method, body string) *httptest.ResponseRecorder {
	t.Helper()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(method, "/calculate", strings.NewReader(body))
	handler.ServeHTTP(recorder, req)
	return recorder
}
func TestBinaryValidation(t *testing.T) {
	handler := BinaryHandler(Add)
	for _, tc := range []struct{ name, body string }{
		{"malformed", `{`}, {"missing", `{"valor1":1}`}, {"extra", `{"valor1":1,"valor2":2,"extra":3}`}, {"null", `{"valor1":null,"valor2":2}`}, {"multiple", `{"valor1":1,"valor2":2}{}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := request(t, handler, "POST", tc.body); got.Code != 400 {
				t.Fatalf("got %d: %s", got.Code, got.Body.String())
			}
		})
	}
	if got := request(t, handler, "GET", `{}`); got.Code != 405 {
		t.Fatalf("got %d", got.Code)
	}
}
