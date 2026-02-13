package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cast"
)

func ScanfInt(msg ...string) int {
	result, _ := pterm.DefaultInteractiveTextInput.Show(msg...)
	return cast.ToInt(result)
}

func Confirm(msg ...string) bool {
	b, _ := pterm.DefaultInteractiveConfirm.Show(msg...)
	return b
}

// SanitizeFilename replaces common unsafe characters with underscores
func SanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		"?", "_",
		"&", "_",
		"=", "_",
		"#", "_",
		":", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"*", "_",
		"|", "_",
	)

	return replacer.Replace(name)
}

func RunCommand(name string, args []string) (string, error) {
	cmd := exec.Command(name, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %w, stderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}
