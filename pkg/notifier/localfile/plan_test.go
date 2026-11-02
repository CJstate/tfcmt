package localfile

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/notifier"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
)

func TestNotifyPlanWritesFile(t *testing.T) {
	t.Parallel()
	outputFile := filepath.Join(t.TempDir(), "plan.md")
	client, err := NewClient(&Config{
		OutputFile:         outputFile,
		Parser:             terraform.NewPlanParser(),
		Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
		ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanParseErrorTemplate),
		DisableLabel:       true,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Notify.Plan(context.Background(), slog.New(slog.DiscardHandler), &notifier.ParamExec{
		CombinedOutput: "Plan: 1 to add",
		ExitCode:       2,
	}); err != nil {
		t.Fatal(err)
	}

	body, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "Plan: 1 to add") {
		t.Errorf("output file doesn't contain the plan result: %s", body)
	}
}

func TestNotifyApplyWritesFile(t *testing.T) {
	t.Parallel()
	outputFile := filepath.Join(t.TempDir(), "apply.md")
	client, err := NewClient(&Config{
		OutputFile:         outputFile,
		Parser:             terraform.NewApplyParser(),
		Template:           terraform.NewApplyTemplate(terraform.DefaultApplyTemplate),
		ParseErrorTemplate: terraform.NewApplyParseErrorTemplate(terraform.DefaultApplyParseErrorTemplate),
		DisableLabel:       true,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := client.Notify.Apply(context.Background(), slog.New(slog.DiscardHandler), &notifier.ParamExec{
		CombinedOutput: "Apply complete! Resources: 1 added, 0 changed, 0 destroyed.",
		ExitCode:       0,
	}); err != nil {
		t.Fatal(err)
	}

	body, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "Apply complete!") {
		t.Errorf("output file doesn't contain the apply result: %s", body)
	}
}
