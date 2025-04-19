package events

import "context"

type Producer interface {
	ProducerInconsistentEvent(ctx context.Context, evt InconsistentEvent) error
}

type InconsistentEvent struct {
	ID int64
	// src 以原表为准，dst 以目标表为准
	Direction string
	Type      string
}

const (
	InconsistentEventTypeTargetMissing = "target_missing"
	InconsistentEventTypeBaseMissing   = "base_missing"
	InconsistentEventTypeNEQ           = "neq"
)
