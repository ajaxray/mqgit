package util

import (
	"strings"
)

// LastGitMessage returns trimmed latest git message (if in a git repo)
func LastGitMessage() string {
	if message, err := RunCommand("git log -1 --pretty=%B"); err == nil {
		return strings.TrimSpace(string(message))
	}

	return ""
}

// CurrentGitHash returns short hash of current git commit
func CurrentGitHash() string {
	if hash, err := RunCommand("git rev-parse --short HEAD"); err == nil {
		return strings.TrimSpace(string(hash))
	}

	return ""
}
