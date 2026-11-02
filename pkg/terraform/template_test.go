package terraform

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDefaultUpdatedResourcesTemplate(t *testing.T) {
	t.Parallel()
	wrappedList := func(resources ...string) string {
		return "\n```hcl\n" + strings.Join(resources, "\n") + "\n```\n"
	}
	want := "\n* Create\n" + wrappedList("aws_instance.foo", "aws_instance.bar") +
		"\n* Update\n" + wrappedList("aws_instance.baz") +
		"\n* Delete\n" + wrappedList("aws_instance.qux") +
		"\n* Replace\n" + wrappedList("aws_instance.quux") +
		"\n* Import\n" + wrappedList("aws_instance.corge") +
		"\n* Move\n" + wrappedList("aws_instance.grault => aws_instance.garply")

	for _, testCase := range []struct {
		name         string
		useRawOutput bool
	}{
		{name: "html template", useRawOutput: false},
		{name: "raw output", useRawOutput: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			tpl := &Template{
				CommonTemplate: CommonTemplate{
					UseRawOutput:      testCase.useRawOutput,
					CreatedResources:  []string{"aws_instance.foo", "aws_instance.bar"},
					UpdatedResources:  []string{"aws_instance.baz"},
					DeletedResources:  []string{"aws_instance.qux"},
					ReplacedResources: []string{"aws_instance.quux"},
					ImportedResources: []string{"aws_instance.corge"},
					MovedResources: []*MovedResource{
						{Before: "aws_instance.grault", After: "aws_instance.garply"},
					},
				},
				Template: `{{template "updated_resources" .}}`,
			}

			got, err := tpl.Execute()
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("updated_resources mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDefaultUpdatedResourcesTemplateTruncatesLongList(t *testing.T) {
	t.Parallel()
	const omittedMarker = "OMITTED-MIDDLE"
	tpl := &Template{
		CommonTemplate: CommonTemplate{
			UpdatedResources: []string{
				"START-" + strings.Repeat("a", 40000),
				omittedMarker,
				strings.Repeat("z", 30000) + "-END",
			},
		},
		Template: `{{template "updated_resources" .}}`,
	}

	got, err := tpl.Execute()
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"* Update",
		":warning: **The content is omitted as it is too long.** :warning:",
		"START-",
		"-END",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("updated_resources output doesn't contain %q", want)
		}
	}
	if strings.Contains(got, omittedMarker) {
		t.Errorf("updated_resources output contains truncated content %q", omittedMarker)
	}
}
