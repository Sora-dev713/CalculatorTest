package main

import (
	"calculator-app/internal/orchestrator"
	"log"
	"net/http"
	"os"
	"time"
)

func env(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
func main() {
	timeout, err := time.ParseDuration(env("SERVICE_TIMEOUT", "3s"))
	if err != nil {
		log.Fatal("invalid SERVICE_TIMEOUT: ", err)
	}
	client := &orchestrator.Client{HTTP: &http.Client{Timeout: timeout}, URLs: map[string]string{
		"addition": env("ADDITION_URL", "http://localhost:8081"), "subtraction": env("SUBTRACTION_URL", "http://localhost:8082"),
		"multiplication": env("MULTIPLICATION_URL", "http://localhost:8083"), "division": env("DIVISION_URL", "http://localhost:8084"),
		"power": env("POWER_URL", "http://localhost:8085"), "sqrt": env("SQRT_URL", "http://localhost:8086"), "percent": env("PERCENTAGE_URL", "http://localhost:8087"),
	}}
	port := env("PORT", "8080")
	mux := http.NewServeMux()
	mux.Handle("/api/calculate", orchestrator.Handler(client))
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	srv := &http.Server{Addr: ":" + port, Handler: mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	log.Printf("orchestrator listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}
