package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"calculator-app/internal/operation"
)

func TestSquareRootService(t *testing.T) {
	tests := []struct {
		name, body, expected string
		status               int
	}{
		{"valid square root", `{"valor1":16}`, `"resultado":4`, 200},
		{"negative radicand", `{"valor1":-1}`, `"Status":"ERROR"`, 400},
		{"rejects valor2", `{"valor1":4,"valor2":2}`, `"Status":"ERROR"`, 400},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/calculate", strings.NewReader(test.body))
			operation.UnaryHandler(operation.SquareRoot).ServeHTTP(recorder, request)
			if recorder.Code != test.status || !strings.Contains(recorder.Body.String(), test.expected) {
				t.Fatalf("got %d %s", recorder.Code, recorder.Body.String())
			}
		})
	}
}
