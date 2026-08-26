package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"calculator-app/internal/operation"
)

func TestSubtractionService(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/calculate", strings.NewReader(`{"valor1":-2,"valor2":3}`))
	operation.BinaryHandler(operation.Subtract).ServeHTTP(recorder, request)
	if recorder.Code != 200 || !strings.Contains(recorder.Body.String(), `"resultado":-5`) {
		t.Fatalf("got %d %s", recorder.Code, recorder.Body.String())
	}
}
