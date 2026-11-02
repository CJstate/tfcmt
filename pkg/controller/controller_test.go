package controller

import (
	"context"
	"testing"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/config"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/notifier"
	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/notifier/localfile"
)

func TestGetPlanNotifierOutputComment(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "xxx")
	testCases := []struct {
		name          string
		output        string
		outputComment bool
		disableLabel  bool
		wantMulti     bool
	}{
		{
			name:         "output only",
			output:       "plan.md",
			disableLabel: true,
		},
		{
			name:          "output and comment",
			output:        "plan.md",
			outputComment: true,
			disableLabel:  true,
			wantMulti:     true,
		},
		{
			name:          "output and comment with labels",
			output:        "plan.md",
			outputComment: true,
			wantMulti:     true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := &Controller{
				Config: config.Config{
					Output:        testCase.output,
					OutputComment: testCase.outputComment,
					Terraform: config.Terraform{
						Plan: config.Plan{
							DisableLabel: testCase.disableLabel,
						},
					},
				},
			}
			ntf, err := c.getPlanNotifier(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			_, isMulti := ntf.(*notifier.MultiNotifier)
			if isMulti != testCase.wantMulti {
				t.Errorf("got multi notifier %v but want %v", isMulti, testCase.wantMulti)
			}
			if !testCase.wantMulti {
				if _, ok := ntf.(*localfile.NotifyService); !ok {
					t.Errorf("got notifier %T but want *localfile.NotifyService", ntf)
				}
			}
		})
	}
}

func TestGetPlanNotifierOutputOnlyWithoutToken(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("TFCMT_GITHUB_TOKEN", "")
	c := &Controller{
		Config: config.Config{
			Output: "plan.md",
			Terraform: config.Terraform{
				Plan: config.Plan{
					DisableLabel: true,
				},
			},
		},
	}
	ntf, err := c.getPlanNotifier(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := ntf.(*localfile.NotifyService); !ok {
		t.Errorf("got notifier %T but want *localfile.NotifyService", ntf)
	}
}

func TestGetApplyNotifierOutputOnlyWithoutToken(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("TFCMT_GITHUB_TOKEN", "")
	c := &Controller{
		Config: config.Config{
			Output: "apply.md",
		},
	}
	ntf, err := c.getApplyNotifier(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := ntf.(*localfile.NotifyService); !ok {
		t.Errorf("got notifier %T but want *localfile.NotifyService", ntf)
	}
}

func TestGetApplyNotifierOutputComment(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "xxx")
	testCases := []struct {
		name          string
		output        string
		outputComment bool
		wantMulti     bool
	}{
		{
			name:   "output only",
			output: "apply.md",
		},
		{
			name:          "output and comment",
			output:        "apply.md",
			outputComment: true,
			wantMulti:     true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := &Controller{
				Config: config.Config{
					Output:        testCase.output,
					OutputComment: testCase.outputComment,
				},
			}
			ntf, err := c.getApplyNotifier(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			_, isMulti := ntf.(*notifier.MultiNotifier)
			if isMulti != testCase.wantMulti {
				t.Errorf("got multi notifier %v but want %v", isMulti, testCase.wantMulti)
			}
			if !testCase.wantMulti {
				if _, ok := ntf.(*localfile.NotifyService); !ok {
					t.Errorf("got notifier %T but want *localfile.NotifyService", ntf)
				}
			}
		})
	}
}
