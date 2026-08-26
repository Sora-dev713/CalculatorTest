package orchestrator

import (
	"net/http"

	"calculator-app/internal/api"
	"calculator-app/internal/expression"
)

// CalculateRequest is the public request accepted by the orchestrator.
type CalculateRequest struct {
	Expression *string `json:"expression"`
}

// Handler validates and parses an expression before starting distributed evaluation.
func Handler(client *Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			api.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var request CalculateRequest
		if err := api.DecodeStrict(r, &request); err != nil {
			api.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		if request.Expression == nil {
			api.WriteError(w, http.StatusBadRequest, "expression is required")
			return
		}
		node, err := expression.Parse(*request.Expression)
		if err != nil {
			api.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		result, err := client.Evaluate(r.Context(), node)
		if err != nil {
			api.WriteError(w, http.StatusBadGateway, err.Error())
			return
		}
		api.WriteResult(w, http.StatusOK, result)
	})
}
