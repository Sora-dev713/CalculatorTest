package operation

import (
	"errors"
	"math"
	"net/http"

	"calculator-app/internal/api"
)

// BinaryFunc is the common signature implemented by two-operand calculations.
type BinaryFunc func(float64, float64) (float64, error)

// UnaryFunc is the common signature implemented by one-operand calculations.
type UnaryFunc func(float64) (float64, error)

// BinaryHandler adapts a BinaryFunc to the shared strict HTTP contract.
func BinaryHandler(calculate BinaryFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			api.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var request api.BinaryRequest
		if err := api.DecodeStrict(r, &request); err != nil {
			api.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		if request.Valor1 == nil || request.Valor2 == nil {
			api.WriteError(w, http.StatusBadRequest, "valor1 and valor2 are required")
			return
		}
		if !api.ValidNumber(*request.Valor1) || !api.ValidNumber(*request.Valor2) {
			api.WriteError(w, http.StatusBadRequest, "values must be finite numbers")
			return
		}
		result, err := calculate(*request.Valor1, *request.Valor2)
		if err != nil {
			api.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		if !api.ValidNumber(result) {
			api.WriteError(w, http.StatusBadRequest, "result is not a finite number")
			return
		}
		api.WriteResult(w, http.StatusOK, result)
	})
}

// UnaryHandler adapts a UnaryFunc to the shared strict HTTP contract.
func UnaryHandler(calculate UnaryFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			api.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var request api.UnaryRequest
		if err := api.DecodeStrict(r, &request); err != nil {
			api.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		if request.Valor1 == nil {
			api.WriteError(w, http.StatusBadRequest, "valor1 is required")
			return
		}
		if !api.ValidNumber(*request.Valor1) {
			api.WriteError(w, http.StatusBadRequest, "valor1 must be a finite number")
			return
		}
		result, err := calculate(*request.Valor1)
		if err != nil {
			api.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		if math.IsNaN(result) || math.IsInf(result, 0) {
			api.WriteError(w, http.StatusBadRequest, "result is not a finite number")
			return
		}
		api.WriteResult(w, http.StatusOK, result)
	})
}

// Add returns the sum of a and b.
func Add(a, b float64) (float64, error) { return a + b, nil }

// Subtract returns b subtracted from a.
func Subtract(a, b float64) (float64, error) { return a - b, nil }

// Multiply returns the product of a and b.
func Multiply(a, b float64) (float64, error) { return a * b, nil }

// Divide returns a divided by b and rejects a zero divisor.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero is not allowed")
	}
	return a / b, nil
}

// Power raises a to b and rejects results outside the finite float64 domain.
func Power(a, b float64) (float64, error) {
	result := math.Pow(a, b)
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, errors.New("power is outside the supported numeric domain")
	}
	return result, nil
}

// SquareRoot returns the principal square root and rejects negative values.
func SquareRoot(value float64) (float64, error) {
	if value < 0 {
		return 0, errors.New("square root of a negative number is not allowed")
	}
	return math.Sqrt(value), nil
}

// Percentage returns percentage percent of value; negative percentages are valid.
func Percentage(value, percentage float64) (float64, error) { return value * percentage / 100, nil }
