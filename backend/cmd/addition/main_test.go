package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"calculator-app/internal/operation"
)

func TestAdditionService(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/calculate", strings.NewReader(`{"valor1":1.5,"valor2":2.25}`))
	operation.BinaryHandler(operation.Add).ServeHTTP(recorder, request)
	if recorder.Code != 200 || !strings.Contains(recorder.Body.String(), `"resultado":3.75`) {
		t.Fatalf("got %d %s", recorder.Code, recorder.Body.String())
	}
}
