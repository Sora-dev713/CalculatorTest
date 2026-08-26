package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"calculator-app/internal/operation"
)

func TestDivisionService(t *testing.T) {
	tests := []struct {
		name, body, expected string
		status               int
	}{
		{"valid division", `{"valor1":9,"valor2":3}`, `"resultado":3`, 200},
		{"division by zero", `{"valor1":9,"valor2":0}`, `"Status":"ERROR"`, 400},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/calculate", strings.NewReader(test.body))
			operation.BinaryHandler(operation.Divide).ServeHTTP(recorder, request)
			if recorder.Code != test.status || !strings.Contains(recorder.Body.String(), test.expected) {
				t.Fatalf("got %d %s", recorder.Code, recorder.Body.String())
			}
		})
	}
}
