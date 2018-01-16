package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Prompt receives string input with space on command prompt
func Prompt(question string, defaultAnswer string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(question)
	text, _ := reader.ReadString('\n')
	if text == "\n" {
		return defaultAnswer
	}

	return strings.TrimSuffix(text, "\n")
}

// MakeShellCmd creates a exec.Cmd from string
// Then you can use it for example: output, err := cmd.Output()
func MakeShellCmd(cmd string) *exec.Cmd {
	return exec.Command("sh", "-c", cmd)
}

// RunCommand runs a string as shell command
func RunCommand(cmd string) ([]byte, error) {
	execCmd := MakeShellCmd(cmd)
	return execCmd.Output()
}
