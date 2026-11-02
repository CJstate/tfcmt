package jobsummary

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestNewWriter(t *testing.T) {
	t.Parallel()
	if _, err := NewWriter(""); err == nil {
		t.Error("NewWriter with an empty path should return an error")
	}
	if _, err := NewWriter("summary.md"); err != nil {
		t.Errorf("NewWriter with a path should not return an error: %v", err)
	}
}

func TestNewWriterFromEnv(t *testing.T) {
	path := filepath.Join(t.TempDir(), "summary.md")
	t.Setenv(EnvVar, path)
	writer, err := NewWriterFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if writer.path != path {
		t.Errorf("got %q but want %q", writer.path, path)
	}
}

func TestNewWriterFromEnvMissing(t *testing.T) {
	t.Setenv(EnvVar, "")
	if _, err := NewWriterFromEnv(); err == nil {
		t.Error("NewWriterFromEnv without GITHUB_STEP_SUMMARY should return an error")
	}
}

func TestWriterWrite(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "summary.md")
	writer, err := NewWriter(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := writer.Write("hello"); err != nil {
		t.Fatal(err)
	}
	if err := writer.Write("world"); err != nil {
		t.Fatal(err)
	}
	actual, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(actual), "hello\nworld\n"; got != want {
		t.Errorf("got %q but want %q", got, want)
	}
}

func TestWriterTruncates(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "summary.md")
	writer, err := NewWriter(path)
	if err != nil {
		t.Fatal(err)
	}
	body := strings.Repeat("a", MaxBytes+1024)
	if err := writer.Write(body); err != nil {
		t.Fatal(err)
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if fileInfo.Size() > MaxBytes {
		t.Errorf("the job summary size is %d bytes, which exceeds the limit of %d bytes", fileInfo.Size(), MaxBytes)
	}
	actual, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(string(actual), truncationNotice) {
		t.Errorf("the job summary should end with the truncation notice: %q", truncationNotice)
	}
}

func TestWriterTruncatesWithExistingContent(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "summary.md")
	existing := strings.Repeat("b", MaxBytes-200)
	if err := os.WriteFile(path, []byte(existing), 0o600); err != nil {
		t.Fatal(err)
	}
	writer, err := NewWriter(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := writer.Write(strings.Repeat("a", 1024)); err != nil {
		t.Fatal(err)
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if fileInfo.Size() > MaxBytes {
		t.Errorf("the job summary size is %d bytes, which exceeds the limit of %d bytes", fileInfo.Size(), MaxBytes)
	}
	actual, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(actual), existing) {
		t.Error("the existing job summary content should be preserved")
	}
}

func TestWriterFullSummary(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "summary.md")
	existing := strings.Repeat("b", MaxBytes)
	if err := os.WriteFile(path, []byte(existing), 0o600); err != nil {
		t.Fatal(err)
	}
	writer, err := NewWriter(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := writer.Write("hello"); err != nil {
		t.Fatal(err)
	}
	actual, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(actual) != existing {
		t.Error("a full job summary should not be modified")
	}
}

func TestTruncateUTF8(t *testing.T) {
	t.Parallel()
	body := strings.Repeat("あ", 1000)
	actual := truncateUTF8(body, 101)
	if len(actual) > 101 {
		t.Errorf("the truncated body is %d bytes, which exceeds the limit of 101 bytes", len(actual))
	}
	if !utf8.ValidString(actual) {
		t.Error("the truncated body should be valid UTF-8")
	}
}
