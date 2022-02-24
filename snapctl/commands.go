package snapctl

import (
	"fmt"
	"os/exec"
	"strings"
)

func run(subcommand string, subargs ...string) (string, error) {
	args := []string{subcommand}
	args = append(args, subargs...)

	if debug {
		fmt.Printf("[debug] snapctl %s\n",
			strings.Join(args, " "))
	}

	output, err := exec.Command("snapctl", args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, output)
	}

	return strings.TrimSpace(string(output)), nil
}
