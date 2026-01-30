package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// This file is deprecated. Use cmd/server/main.go instead.
	// For backwards compatibility, we'll delegate to the server binary.

	// If running as a worker, delegate to cmd/worker
	if len(os.Args) > 1 && os.Args[1] == "worker" {
		cmd := exec.Command("go", append([]string{"run", "./cmd/worker"}, os.Args[2:]...)...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running worker: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Default: run server
	cmd := exec.Command("go", append([]string{"run", "./cmd/server"}, os.Args[1:]...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running server: %v\n", err)
		os.Exit(1)
	}
}
