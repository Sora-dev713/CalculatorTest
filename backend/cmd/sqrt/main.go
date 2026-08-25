package main

import (
	"calculator-app/internal/operation"
	"calculator-app/internal/server"
)

func main() { server.Run(operation.UnaryHandler(operation.SquareRoot)) }
