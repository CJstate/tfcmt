package notifier

import (
	"context"
	"log/slog"
)

// MultiNotifier notifies with multiple notifiers in order.
// It is used to output the result to a file and post a comment at the same time.
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier returns a notifier which calls the given notifiers in order.
func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{
		notifiers: notifiers,
	}
}

// Apply notifies the apply result with all notifiers in order.
func (m *MultiNotifier) Apply(ctx context.Context, logger *slog.Logger, param *ParamExec) error {
	for _, n := range m.notifiers {
		if err := n.Apply(ctx, logger, param); err != nil {
			return err
		}
	}
	return nil
}

// Plan notifies the plan result with all notifiers in order.
func (m *MultiNotifier) Plan(ctx context.Context, logger *slog.Logger, param *ParamExec) error {
	for _, n := range m.notifiers {
		if err := n.Plan(ctx, logger, param); err != nil {
			return err
		}
	}
	return nil
}
