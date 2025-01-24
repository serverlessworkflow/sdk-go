package impl

import (
	"context"
	"fmt"
	"github.com/serverlessworkflow/sdk-go/v3/model"
)

type StatusPhase string

const (
	PendingStatus   StatusPhase = "pending"
	RunningStatus   StatusPhase = "running"
	WaitingStatus   StatusPhase = "waiting"
	CancelledStatus StatusPhase = "cancelled"
	FaultedStatus   StatusPhase = "faulted"
	CompletedStatus StatusPhase = "completed"
)

var _ WorkflowRunner = &workflowRunnerImpl{}

type WorkflowRunner interface {
	GetWorkflow() *model.Workflow
	Run(input map[string]interface{}) (output map[string]interface{}, err error)
}

func NewDefaultRunner(workflow *model.Workflow) WorkflowRunner {
	// later we can implement the opts pattern to define context timeout, deadline, cancel, etc.
	// also fetch from the workflow model this information
	ctx := WithExecutorContext(context.Background(), &ExecutorContext{})
	return &workflowRunnerImpl{
		Workflow: workflow,
		Context:  ctx,
	}
}

type workflowRunnerImpl struct {
	Workflow *model.Workflow
	Context  context.Context
}

func (wr *workflowRunnerImpl) GetWorkflow() *model.Workflow {
	return wr.Workflow
}

// Run the workflow.
// TODO: Sync execution, we think about async later
func (wr *workflowRunnerImpl) Run(input map[string]interface{}) (output map[string]interface{}, err error) {
	output = make(map[string]interface{})
	if input == nil {
		input = make(map[string]interface{})
	}

	// TODO: validates input via wr.Workflow.Input.Schema

	wfCtx, err := GetExecutorContext(wr.Context)
	if err != nil {
		return nil, err
	}
	wfCtx.SetInput(input)
	wfCtx.SetOutput(output)

	// TODO: process wr.Workflow.Input.From, the result we set to WorkFlowCtx
	wfCtx.SetWorkflowCtx(input)

	// Run tasks
	// For each task, execute.
	if wr.Workflow.Do != nil {
		for _, taskItem := range *wr.Workflow.Do {
			switch task := taskItem.Task.(type) {
			case *model.SetTask:
				exec, err := NewSetTaskExecutor(taskItem.Key, task)
				if err != nil {
					return nil, err
				}
				output, err = exec.Exec(wfCtx.GetWorkflowCtx())
				if err != nil {
					return nil, err
				}
				wfCtx.SetWorkflowCtx(output)
			default:
				return nil, fmt.Errorf("workflow does not support task '%T' named '%s'", task, taskItem.Key)
			}
		}
	}

	// Process output and return

	return output, err
}
