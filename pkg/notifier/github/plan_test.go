package github

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/go-github/v93/github"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/notifier"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
)

func TestNotifyPlanSkipsLabelUpdate(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "xxx")
	cfg := Config{
		Owner: "owner",
		Repo:  "repo",
		PR: PullRequest{
			Revision: "revision",
			Number:   1,
		},
		Parser:             terraform.NewPlanParser(),
		Template:           terraform.NewPlanTemplate(terraform.DefaultPlanTemplate),
		ParseErrorTemplate: terraform.NewPlanParseErrorTemplate(terraform.DefaultPlanParseErrorTemplate),
		ResultLabels: ResultLabels{
			AddOrUpdateLabel: "add-or-update",
		},
		SkipLabelUpdate: true,
	}
	client, err := NewClient(t.Context(), &cfg)
	if err != nil {
		t.Fatal(err)
	}
	api := &fakeAPI{
		FakeIssuesCreateComment: func(_ context.Context, _ int, _ github.IssueCommentRequest) (*github.IssueComment, *github.Response, error) {
			return &github.IssueComment{}, nil, nil
		},
	}
	client.API = api
	if err := client.Notify.Plan(t.Context(), slog.New(slog.DiscardHandler), &notifier.ParamExec{
		CombinedOutput: "Plan: 1 to add",
		ExitCode:       2,
	}); err != nil {
		t.Fatal(err)
	}
}
