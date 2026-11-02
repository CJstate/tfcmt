// Package jobsummary writes tfcmt results to the GitHub Actions job summary.
package jobsummary

import (
	"errors"
	"fmt"
	"os"
	"unicode/utf8"
)

// EnvVar is the environment variable GitHub Actions sets to the path of the job summary file.
// https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#adding-a-job-summary
const EnvVar = "GITHUB_STEP_SUMMARY"

// MaxBytes is the maximum size of a GitHub Actions job summary (1 MiB).
const MaxBytes = 1 << 20

const (
	filePermission = os.FileMode(0o644)
	// truncationNotice is appended when a result is too large to fit in the job summary.
	truncationNotice = "\n\n> Output truncated to fit the GitHub Actions job summary limit (1 MiB).\n"
)

// Writer appends rendered results to a GitHub Actions job summary file.
type Writer struct {
	path string
}

// NewWriter returns a Writer for the given job summary file path.
func NewWriter(path string) (*Writer, error) {
	if path == "" {
		return nil, errors.New("the path of the GitHub Actions job summary is empty")
	}
	return &Writer{path: path}, nil
}

// NewWriterFromEnv returns a Writer for the file specified by GITHUB_STEP_SUMMARY.
func NewWriterFromEnv() (*Writer, error) {
	return NewWriter(os.Getenv(EnvVar))
}

// Write appends the body to the job summary, truncating it if the result would exceed
// the GitHub Actions 1 MiB job summary limit.
func (w *Writer) Write(body string) error {
	data := body + "\n"
	available := MaxBytes - w.currentSize()
	if available <= 0 {
		return nil
	}
	if len(data) > available {
		data = truncate(data, available)
	}
	if data == "" {
		return nil
	}
	file, err := os.OpenFile(w.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePermission)
	if err != nil {
		return fmt.Errorf("open the GitHub Actions job summary: %w", err)
	}
	defer file.Close()
	if _, err := file.WriteString(data); err != nil {
		return fmt.Errorf("write the GitHub Actions job summary: %w", err)
	}
	return nil
}

func (w *Writer) currentSize() int {
	fileInfo, err := os.Stat(w.path)
	if err != nil {
		return 0
	}
	if fileInfo.Size() > MaxBytes {
		return MaxBytes
	}
	return int(fileInfo.Size())
}

func truncate(body string, limit int) string {
	if limit <= len(truncationNotice) {
		return truncationNotice[:limit]
	}
	return truncateUTF8(body, limit-len(truncationNotice)) + truncationNotice
}

func truncateUTF8(body string, limit int) string {
	if len(body) <= limit {
		return body
	}
	body = body[:limit]
	for !utf8.ValidString(body) {
		body = body[:len(body)-1]
	}
	return body
}
