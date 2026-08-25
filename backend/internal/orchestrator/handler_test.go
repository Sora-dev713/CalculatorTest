package orchestrator

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"calculator-app/internal/expression"
	"calculator-app/internal/operation"
)

func TestEvaluateRealOperationServices(t *testing.T) {
	services := map[string]http.Handler{
		"addition": operation.BinaryHandler(operation.Add), "subtraction": operation.BinaryHandler(operation.Subtract),
		"multiplication": operation.BinaryHandler(operation.Multiply), "division": operation.BinaryHandler(operation.Divide),
		"power": operation.BinaryHandler(operation.Power), "sqrt": operation.UnaryHandler(operation.SquareRoot),
		"percent": operation.BinaryHandler(operation.Percentage),
	}
	urls := map[string]string{}
	servers := []*httptest.Server{}
	for name, handler := range services {
		server := httptest.NewServer(handler)
		servers = append(servers, server)
		urls[name] = server.URL
	}
	defer func() {
		for _, server := range servers {
			server.Close()
		}
	}()
	client := &Client{HTTP: http.DefaultClient, URLs: urls}
	for input, expected := range map[string]float64{
		"2+3*4": 14, "(2+3)*4": 20, "2^3^2": 512, "sqrt(16)+2": 6, "percent(200,10)": 20, "-2^2": -4,
	} {
		t.Run(input, func(t *testing.T) {
			node, err := expression.Parse(input)
			if err != nil {
				t.Fatal(err)
			}
			result, err := client.Evaluate(t.Context(), node)
			if err != nil {
				t.Fatal(err)
			}
			if result != expected {
				t.Fatalf("got %v, want %v", result, expected)
			}
		})
	}
}

func mockOperations(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/calculate" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.Header.Get("X-Test") {
		default:
			_, _ = w.Write([]byte(`{"Status":"ok","resultado":42}`))
		}
	}))
}
func TestHandlerEvaluatesThroughService(t *testing.T) {
	server := mockOperations(t)
	defer server.Close()
	urls := map[string]string{"addition": server.URL, "subtraction": server.URL, "multiplication": server.URL, "division": server.URL, "power": server.URL, "sqrt": server.URL, "percent": server.URL}
	handler := Handler(&Client{HTTP: server.Client(), URLs: urls})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest("POST", "/api/calculate", strings.NewReader(`{"expression":"2+3"}`)))
	if recorder.Code != 200 || !strings.Contains(recorder.Body.String(), `"resultado":42`) {
		t.Fatalf("got %d %s", recorder.Code, recorder.Body.String())
	}
}
func TestHandlerValidatesBeforeCallingServices(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls++ }))
	defer server.Close()
	handler := Handler(&Client{HTTP: server.Client(), URLs: map[string]string{"addition": server.URL}})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest("POST", "/api/calculate", strings.NewReader(`{"expression":"2+bad"}`)))
	if recorder.Code != 400 || calls != 0 {
		t.Fatalf("got status %d and %d calls", recorder.Code, calls)
	}
}
func TestHandlerPropagatesServiceError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"Status":"ERROR","Error":"division by zero is not allowed"}`))
	}))
	defer server.Close()
	handler := Handler(&Client{HTTP: server.Client(), URLs: map[string]string{"division": server.URL}})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest("POST", "/api/calculate", strings.NewReader(`{"expression":"1/0"}`)))
	if recorder.Code != 502 || !strings.Contains(recorder.Body.String(), "division by zero") {
		t.Fatalf("got %d %s", recorder.Code, recorder.Body.String())
	}
}
func TestClientTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(50 * time.Millisecond)
		_, _ = w.Write([]byte(`{"Status":"ok","resultado":1}`))
	}))
	defer server.Close()
	client := &Client{HTTP: &http.Client{Timeout: time.Millisecond}, URLs: map[string]string{"addition": server.URL}}
	handler := Handler(client)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest("POST", "/api/calculate", strings.NewReader(`{"expression":"1+1"}`)))
	if recorder.Code != 502 || !strings.Contains(recorder.Body.String(), "unavailable") {
		t.Fatalf("got %d %s", recorder.Code, recorder.Body.String())
	}
}
