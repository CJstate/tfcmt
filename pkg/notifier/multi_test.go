package notifier

import (
	"context"
	"errors"
	"log/slog"
	"testing"
)

type fakeNotifier struct {
	planCalls  int
	applyCalls int
	planErr    error
}

func (f *fakeNotifier) Plan(_ context.Context, _ *slog.Logger, _ *ParamExec) error {
	f.planCalls++
	return f.planErr
}

func (f *fakeNotifier) Apply(_ context.Context, _ *slog.Logger, _ *ParamExec) error {
	f.applyCalls++
	return nil
}

func TestMultiNotifierPlan(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.DiscardHandler)
	param := &ParamExec{}

	t.Run("calls all notifiers in order", func(t *testing.T) {
		t.Parallel()
		first := &fakeNotifier{}
		second := &fakeNotifier{}
		multi := NewMultiNotifier(first, second)
		if err := multi.Plan(t.Context(), logger, param); err != nil {
			t.Fatal(err)
		}
		if first.planCalls != 1 || second.planCalls != 1 {
			t.Errorf("got calls %d and %d, want 1 and 1", first.planCalls, second.planCalls)
		}
		if first.applyCalls != 0 || second.applyCalls != 0 {
			t.Errorf("unexpected apply calls: %d and %d", first.applyCalls, second.applyCalls)
		}
	})

	t.Run("stops when a notifier returns an error", func(t *testing.T) {
		t.Parallel()
		wantErr := errors.New("failed to post a comment")
		first := &fakeNotifier{}
		second := &fakeNotifier{planErr: wantErr}
		third := &fakeNotifier{}
		multi := NewMultiNotifier(first, second, third)
		if err := multi.Plan(t.Context(), logger, param); !errors.Is(err, wantErr) {
			t.Errorf("got error %v, want %v", err, wantErr)
		}
		if first.planCalls != 1 || second.planCalls != 1 || third.planCalls != 0 {
			t.Errorf("got calls %d, %d and %d, want 1, 1 and 0", first.planCalls, second.planCalls, third.planCalls)
		}
	})
}

func TestMultiNotifierApply(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.DiscardHandler)
	param := &ParamExec{}

	first := &fakeNotifier{}
	second := &fakeNotifier{}
	multi := NewMultiNotifier(first, second)
	if err := multi.Apply(t.Context(), logger, param); err != nil {
		t.Fatal(err)
	}
	if first.applyCalls != 1 || second.applyCalls != 1 {
		t.Errorf("got calls %d and %d, want 1 and 1", first.applyCalls, second.applyCalls)
	}
	if first.planCalls != 0 || second.planCalls != 0 {
		t.Errorf("unexpected plan calls: %d and %d", first.planCalls, second.planCalls)
	}
}
