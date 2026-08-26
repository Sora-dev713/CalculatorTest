package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"calculator-app/internal/operation"
)

func TestPowerService(t *testing.T) {
	tests := []struct {
		name, body, expected string
		status               int
	}{
		{"valid power", `{"valor1":2,"valor2":8}`, `"resultado":256`, 200},
		{"invalid domain", `{"valor1":-1,"valor2":0.5}`, `"Status":"ERROR"`, 400},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/calculate", strings.NewReader(test.body))
			operation.BinaryHandler(operation.Power).ServeHTTP(recorder, request)
			if recorder.Code != test.status || !strings.Contains(recorder.Body.String(), test.expected) {
				t.Fatalf("got %d %s", recorder.Code, recorder.Body.String())
			}
		})
	}
}
