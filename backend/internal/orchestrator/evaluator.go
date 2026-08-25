package orchestrator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	"calculator-app/internal/api"
	"calculator-app/internal/expression"
)

type Client struct {
	HTTP *http.Client
	URLs map[string]string
}

func (c *Client) Evaluate(ctx context.Context, node *expression.Node) (float64, error) {
	switch node.Kind {
	case expression.Number:
		return node.Value, nil
	case expression.Unary:
		value, err := c.Evaluate(ctx, node.Right)
		if err != nil {
			return 0, err
		}
		if node.Operator == '-' {
			return -value, nil
		}
		return value, nil
	case expression.Binary:
		left, err := c.Evaluate(ctx, node.Left)
		if err != nil {
			return 0, err
		}
		right, err := c.Evaluate(ctx, node.Right)
		if err != nil {
			return 0, err
		}
		services := map[rune]string{'+': "addition", '-': "subtraction", '*': "multiplication", '/': "division", '^': "power"}
		return c.call(ctx, services[node.Operator], map[string]float64{"valor1": left, "valor2": right})
	case expression.Function:
		first, err := c.Evaluate(ctx, node.Args[0])
		if err != nil {
			return 0, err
		}
		body := map[string]float64{"valor1": first}
		if node.Name == "percent" {
			second, err := c.Evaluate(ctx, node.Args[1])
			if err != nil {
				return 0, err
			}
			body["valor2"] = second
		}
		return c.call(ctx, node.Name, body)
	default:
		return 0, fmt.Errorf("invalid expression node")
	}
}

func (c *Client) call(ctx context.Context, service string, body any) (float64, error) {
	url, ok := c.URLs[service]
	if !ok {
		return 0, fmt.Errorf("service %s is not configured", service)
	}
	payload, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(url, "/")+"/calculate", bytes.NewReader(payload))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := c.HTTP.Do(req)
	if err != nil {
		return 0, fmt.Errorf("%s service is unavailable: %w", service, err)
	}
	defer response.Body.Close()
	data, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return 0, fmt.Errorf("invalid response from %s", service)
	}
	var result api.Response
	if err := json.Unmarshal(data, &result); err != nil {
		return 0, fmt.Errorf("invalid response from %s", service)
	}
	if response.StatusCode != http.StatusOK || result.Status != "ok" || result.Resultado == nil {
		if result.Error != "" {
			return 0, fmt.Errorf("%s", result.Error)
		}
		return 0, fmt.Errorf("%s service returned status %d", service, response.StatusCode)
	}
	if math.IsNaN(*result.Resultado) || math.IsInf(*result.Resultado, 0) {
		return 0, fmt.Errorf("%s returned a non-finite result", service)
	}
	return *result.Resultado, nil
}
