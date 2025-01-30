package impl

import "time"

type StatusPhase string

const (
	// PendingStatus The workflow/task has been initiated and is pending execution.
	PendingStatus StatusPhase = "pending"
	// RunningStatus The workflow/task is currently in progress.
	RunningStatus StatusPhase = "running"
	// WaitingStatus The workflow/task execution is temporarily paused, awaiting either inbound event(s) or a specified time interval as defined by a wait task.
	WaitingStatus StatusPhase = "waiting"
	// SuspendedStatus The workflow/task execution has been manually paused by a user and will remain halted until explicitly resumed.
	SuspendedStatus StatusPhase = "suspended"
	// CancelledStatus The workflow/task execution has been terminated before completion.
	CancelledStatus StatusPhase = "cancelled"
	// FaultedStatus The workflow/task execution has encountered an error.
	FaultedStatus StatusPhase = "faulted"
	// CompletedStatus The workflow/task ran to completion.
	CompletedStatus StatusPhase = "completed"
)

func (s StatusPhase) String() string {
	return string(s)
}

type StatusPhaseLog struct {
	Timestamp int64       `json:"timestamp"`
	Status    StatusPhase `json:"status"`
}

func NewStatusPhaseLog(status StatusPhase) StatusPhaseLog {
	return StatusPhaseLog{
		Status:    status,
		Timestamp: time.Now().UnixMilli(),
	}
}
