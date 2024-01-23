package main

import (
	"github.com/CGPRE-SEPLAN-RR/fiplan-api/internal/server"
	"fmt"
)

func main() {
	server := server.NewServer()

	err := server.ListenAndServe()

	if err != nil {
		panic(fmt.Sprintf("O servidor não pôde ser iniciado: %s", err))
	}
}
