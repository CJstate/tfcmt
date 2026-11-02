package cli

import (
	"testing"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/config"
)

func TestParseOptsPlanJobSummary(t *testing.T) { //nolint:tparallel
	testCases := []struct {
		name string
		args *PlanArgs
		cfg  config.Config
		want bool
	}{
		{
			name: "flag is enabled",
			args: &PlanArgs{GlobalArgs: &GlobalArgs{JobSummary: true, JobSummaryCount: 1}},
			want: true,
		},
		{
			name: "flag is disabled explicitly",
			args: &PlanArgs{GlobalArgs: &GlobalArgs{JobSummary: false, JobSummaryCount: 1}},
			cfg:  config.Config{JobSummary: true},
			want: false,
		},
		{
			name: "flag is not set",
			args: &PlanArgs{GlobalArgs: &GlobalArgs{}},
			cfg:  config.Config{JobSummary: true},
			want: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cfg := testCase.cfg
			if err := parseOptsPlan(testCase.args, &cfg, nil); err != nil {
				t.Fatal(err)
			}
			if cfg.JobSummary != testCase.want {
				t.Errorf("got %v but want %v", cfg.JobSummary, testCase.want)
			}
		})
	}
}

func TestParseOptsApplyJobSummary(t *testing.T) { //nolint:tparallel
	testCases := []struct {
		name string
		args *ApplyArgs
		cfg  config.Config
		want bool
	}{
		{
			name: "flag is enabled",
			args: &ApplyArgs{GlobalArgs: &GlobalArgs{JobSummary: true, JobSummaryCount: 1}},
			want: true,
		},
		{
			name: "flag is disabled explicitly",
			args: &ApplyArgs{GlobalArgs: &GlobalArgs{JobSummary: false, JobSummaryCount: 1}},
			cfg:  config.Config{JobSummary: true},
			want: false,
		},
		{
			name: "flag is not set",
			args: &ApplyArgs{GlobalArgs: &GlobalArgs{}},
			cfg:  config.Config{JobSummary: true},
			want: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cfg := testCase.cfg
			if err := parseOptsApply(testCase.args, &cfg, nil); err != nil {
				t.Fatal(err)
			}
			if cfg.JobSummary != testCase.want {
				t.Errorf("got %v but want %v", cfg.JobSummary, testCase.want)
			}
		})
	}
}
