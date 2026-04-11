package examplesrunner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunAll executes every example generator main file from the repository root.
func RunAll() error {
	commands := [][]string{
		{"go", "run", "./examples/alignment/main.go"},
		{"go", "run", "./examples/flex/flex-test.go"},
		{"go", "run", "./examples/bor-bg-shad/bor-bg-shad.go"},
		{"go", "run", "./examples/tables/table-test.go"},
		{"go", "run", "./examples/phase2/phase2-test.go"},
		{"go", "run", "./examples/fullscale/form2-test.go"},
	}

	for _, cmdArgs := range commands {
		if len(cmdArgs) == 0 {
			continue
		}

		fmt.Printf("Running: %s\n", strings.Join(cmdArgs, " "))
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed command %q: %w", strings.Join(cmdArgs, " "), err)
		}
	}

	return nil
}
