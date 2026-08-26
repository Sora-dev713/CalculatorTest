package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"calculator-app/internal/operation"
)

func TestPercentageService(t *testing.T) {
	tests := []struct{ name, body, expected string }{
		{"positive percentage", `{"valor1":200,"valor2":10}`, `"resultado":20`},
		{"negative percentage", `{"valor1":200,"valor2":-10}`, `"resultado":-20`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/calculate", strings.NewReader(test.body))
			operation.BinaryHandler(operation.Percentage).ServeHTTP(recorder, request)
			if recorder.Code != 200 || !strings.Contains(recorder.Body.String(), test.expected) {
				t.Fatalf("got %d %s", recorder.Code, recorder.Body.String())
			}
		})
	}
}
