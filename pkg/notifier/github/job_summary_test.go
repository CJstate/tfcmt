package github

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/notifier"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/notifier/jobsummary"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
)

func TestNotifyPlanWritesJobSummary(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "xxx")
	path := filepath.Join(t.TempDir(), "summary.md")
	writer, err := jobsummary.NewWriter(path)
	if err != nil {
		t.Fatal(err)
	}
	cfg := newFakeConfig()
	cfg.ParseErrorTemplate = terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanTemplate)
	cfg.JobSummary = writer
	client, err := NewClient(t.Context(), &cfg)
	if err != nil {
		t.Fatal(err)
	}
	api := newFakeAPI()
	client.API = &api
	if err := client.Notify.Plan(t.Context(), slog.New(slog.DiscardHandler), &notifier.ParamExec{
		Stdout:   "Plan: 1 to add",
		ExitCode: 0,
	}); err != nil {
		t.Fatal(err)
	}
	actual, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(actual), "## Plan Result") {
		t.Errorf("the job summary should contain the plan result: %s", actual)
	}
}
