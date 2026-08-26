package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"calculator-app/internal/operation"
)

func TestMultiplicationService(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/calculate", strings.NewReader(`{"valor1":4,"valor2":3}`))
	operation.BinaryHandler(operation.Multiply).ServeHTTP(recorder, request)
	if recorder.Code != 200 || !strings.Contains(recorder.Body.String(), `"resultado":12`) {
		t.Fatalf("got %d %s", recorder.Code, recorder.Body.String())
	}
}
