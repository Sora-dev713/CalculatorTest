package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
)

type BinaryRequest struct {
	Valor1 *float64 `json:"valor1"`
	Valor2 *float64 `json:"valor2"`
}

type UnaryRequest struct {
	Valor1 *float64 `json:"valor1"`
}

type Response struct {
	Status    string   `json:"Status"`
	Resultado *float64 `json:"resultado,omitempty"`
	Error     string   `json:"Error,omitempty"`
}

func WriteResult(w http.ResponseWriter, status int, value float64) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Response{Status: "ok", Resultado: &value})
}

func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Response{Status: "ERROR", Error: message})
}

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

func ValidNumber(value float64) bool { return !math.IsNaN(value) && !math.IsInf(value, 0) }
