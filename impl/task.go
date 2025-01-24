package impl

import (
	"fmt"
	"github.com/serverlessworkflow/sdk-go/v3/expr"
	"github.com/serverlessworkflow/sdk-go/v3/model"
)

var _ TaskExecutor = &SetTaskExecutor{}

type TaskExecutor interface {
	Exec(input map[string]interface{}) (map[string]interface{}, error)
}

type SetTaskExecutor struct {
	Task     *model.SetTask
	TaskName string
}

func NewSetTaskExecutor(taskName string, task *model.SetTask) (*SetTaskExecutor, error) {
	if task == nil || task.Set == nil {
		return nil, fmt.Errorf("no set configuration provided for SetTask %s", taskName)
	}
	return &SetTaskExecutor{
		Task:     task,
		TaskName: taskName,
	}, nil
}

func (s *SetTaskExecutor) Exec(input map[string]interface{}) (output map[string]interface{}, err error) {
	setObject := deepClone(s.Task.Set)
	result, err := expr.TraverseAndEvaluate(setObject, input)
	if err != nil {
		return nil, fmt.Errorf("failed to execute Set task '%s': %w", s.TaskName, err)
	}

	output, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected output to be a map[string]interface{}, but got a different type. Got: %v", result)
	}

	return output, nil
}
