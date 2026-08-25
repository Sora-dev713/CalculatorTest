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
func TestBinaryOperations(t *testing.T) {
	tests := []struct {
		name   string
		fn     BinaryFunc
		body   string
		status int
		result string
	}{
		{"addition decimals", Add, `{"valor1":1.5,"valor2":2.25}`, 200, `"resultado":3.75`},
		{"subtraction negatives", Subtract, `{"valor1":-2,"valor2":3}`, 200, `"resultado":-5`},
		{"multiplication", Multiply, `{"valor1":4,"valor2":3}`, 200, `"resultado":12`},
		{"division", Divide, `{"valor1":9,"valor2":3}`, 200, `"resultado":3`},
		{"division by zero", Divide, `{"valor1":9,"valor2":0}`, 400, `"Status":"ERROR"`},
		{"power", Power, `{"valor1":2,"valor2":8}`, 200, `"resultado":256`},
		{"invalid power", Power, `{"valor1":-1,"valor2":0.5}`, 400, `"Status":"ERROR"`},
		{"percentage negative", Percentage, `{"valor1":200,"valor2":-10}`, 200, `"resultado":-20`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := request(t, BinaryHandler(tc.fn), "POST", tc.body)
			if got.Code != tc.status || !strings.Contains(got.Body.String(), tc.result) {
				t.Fatalf("got %d %s", got.Code, got.Body.String())
			}
		})
	}
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
func TestSquareRoot(t *testing.T) {
	handler := UnaryHandler(SquareRoot)
	if got := request(t, handler, "POST", `{"valor1":16}`); got.Code != 200 || !strings.Contains(got.Body.String(), `"resultado":4`) {
		t.Fatalf("got %d %s", got.Code, got.Body.String())
	}
	if got := request(t, handler, "POST", `{"valor1":-1}`); got.Code != 400 {
		t.Fatalf("got %d", got.Code)
	}
	if got := request(t, handler, "POST", `{"valor1":4,"valor2":2}`); got.Code != 400 {
		t.Fatalf("got %d", got.Code)
	}
}
