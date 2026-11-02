package cli

import (
	"testing"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/config"
)

func TestParseOptsPlanComment(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
		args *PlanArgs
		cfg  config.Config
		want bool
	}{
		{
			name: "flag enables comment",
			args: &PlanArgs{
				GlobalArgs: &GlobalArgs{
					Output:       "plan.md",
					Comment:      true,
					CommentCount: 1,
				},
			},
			want: true,
		},
		{
			name: "config enables comment",
			args: &PlanArgs{
				GlobalArgs: &GlobalArgs{
					Output: "plan.md",
				},
			},
			cfg:  config.Config{OutputComment: true},
			want: true,
		},
		{
			name: "flag explicitly disables comment from config",
			args: &PlanArgs{
				GlobalArgs: &GlobalArgs{
					Output:       "plan.md",
					Comment:      false,
					CommentCount: 1,
				},
			},
			cfg:  config.Config{OutputComment: true},
			want: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			cfg := testCase.cfg
			if err := parseOptsPlan(testCase.args, &cfg, nil); err != nil {
				t.Fatal(err)
			}
			if cfg.OutputComment != testCase.want {
				t.Errorf("got %v but want %v", cfg.OutputComment, testCase.want)
			}
			if cfg.Output != "plan.md" {
				t.Errorf("got output %q but want %q", cfg.Output, "plan.md")
			}
		})
	}
}

func TestParseOptsApplyComment(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
		args *ApplyArgs
		cfg  config.Config
		want bool
	}{
		{
			name: "flag enables comment",
			args: &ApplyArgs{
				GlobalArgs: &GlobalArgs{
					Output:       "apply.md",
					Comment:      true,
					CommentCount: 1,
				},
			},
			want: true,
		},
		{
			name: "config enables comment",
			args: &ApplyArgs{
				GlobalArgs: &GlobalArgs{
					Output: "apply.md",
				},
			},
			cfg:  config.Config{OutputComment: true},
			want: true,
		},
		{
			name: "flag explicitly disables comment from config",
			args: &ApplyArgs{
				GlobalArgs: &GlobalArgs{
					Output:       "apply.md",
					Comment:      false,
					CommentCount: 1,
				},
			},
			cfg:  config.Config{OutputComment: true},
			want: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			cfg := testCase.cfg
			if err := parseOptsApply(testCase.args, &cfg, nil); err != nil {
				t.Fatal(err)
			}
			if cfg.OutputComment != testCase.want {
				t.Errorf("got %v but want %v", cfg.OutputComment, testCase.want)
			}
			if cfg.Output != "apply.md" {
				t.Errorf("got output %q but want %q", cfg.Output, "apply.md")
			}
		})
	}
}
