package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
)

// BinaryRequest is the wire contract shared by all two-operand services.
// Pointers distinguish a missing field from a valid numeric zero.
type BinaryRequest struct {
	Valor1 *float64 `json:"valor1"`
	Valor2 *float64 `json:"valor2"`
}

// UnaryRequest is the wire contract for operations such as square root.
type UnaryRequest struct {
	Valor1 *float64 `json:"valor1"`
}

// Response keeps success and error payloads consistent across all services.
type Response struct {
	Status    string   `json:"Status"`
	Resultado *float64 `json:"resultado,omitempty"`
	Error     string   `json:"Error,omitempty"`
}

// WriteResult serializes a successful calculation using the public JSON contract.
func WriteResult(w http.ResponseWriter, status int, value float64) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Response{Status: "ok", Resultado: &value})
}

// WriteError serializes an error without exposing implementation details.
func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Response{Status: "ERROR", Error: message})
}

// DecodeStrict accepts exactly one JSON object and rejects unknown fields.
func DecodeStrict(r *http.Request, target any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON object")
	}
	return nil
}

// ValidNumber reports whether value can be safely represented in a JSON response.
func ValidNumber(value float64) bool { return !math.IsNaN(value) && !math.IsInf(value, 0) }
