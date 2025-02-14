// Copyright 2025 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
