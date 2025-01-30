package impl

import (
	"fmt"
	"github.com/serverlessworkflow/sdk-go/v3/expr"
	"github.com/serverlessworkflow/sdk-go/v3/model"
)

var _ TaskRunner = &SetTaskRunner{}

type TaskRunner interface {
	Run(input interface{}) (interface{}, error)
	GetTaskName() string
}

func NewSetTaskRunner(taskName string, task *model.SetTask) (*SetTaskRunner, error) {
	if task == nil || task.Set == nil {
		return nil, model.NewErrValidation(fmt.Errorf("no set configuration provided for SetTask %s", taskName), taskName)
	}
	return &SetTaskRunner{
		Task:     task,
		TaskName: taskName,
	}, nil
}

type SetTaskRunner struct {
	Task     *model.SetTask
	TaskName string
}

func (s *SetTaskRunner) GetTaskName() string {
	return s.TaskName
}

func (s *SetTaskRunner) String() string {
	return fmt.Sprintf("SetTaskRunner{Task: %s}", s.GetTaskName())
}

func (s *SetTaskRunner) Run(input interface{}) (output interface{}, err error) {
	setObject := deepClone(s.Task.Set)
	result, err := expr.TraverseAndEvaluate(setObject, input)
	if err != nil {
		return nil, model.NewErrExpression(err, s.TaskName)
	}

	output, ok := result.(map[string]interface{})
	if !ok {
		return nil, model.NewErrRuntime(fmt.Errorf("expected output to be a map[string]interface{}, but got a different type. Got: %v", result), s.TaskName)
	}

	return output, nil
}
