package main

import (
	"log"

	"github.com/Viswesh934/gotei/internal/examplesrunner"
)

func main() {
	if err := examplesrunner.RunAll(); err != nil {
		log.Fatalf("examples runner failed: %v", err)
	}
}
