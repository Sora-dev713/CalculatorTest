package main

import (
	"calculator-app/internal/operation"
	"calculator-app/internal/server"
)

func main() { server.Run(operation.BinaryHandler(operation.Divide)) }
