package impl

import (
	"fmt"
	"github.com/serverlessworkflow/sdk-go/v3/expr"
	"github.com/serverlessworkflow/sdk-go/v3/model"
)

var _ TaskRunner = &SetTaskRunner{}
var _ TaskRunner = &RaiseTaskRunner{}

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

func NewRaiseTaskRunner(taskName string, task *model.RaiseTask) (*RaiseTaskRunner, error) {
	if task == nil || task.Raise.Error.Definition == nil {
		return nil, model.NewErrValidation(fmt.Errorf("no raise configuration provided for RaiseTask %s", taskName), taskName)
	}
	return &RaiseTaskRunner{
		Task:     task,
		TaskName: taskName,
	}, nil
}

type RaiseTaskRunner struct {
	Task     *model.RaiseTask
	TaskName string
}

var raiseErrFuncMapping = map[string]func(error, string) *model.Error{
	model.ErrorTypeAuthentication: model.NewErrAuthentication,
	model.ErrorTypeValidation:     model.NewErrValidation,
	model.ErrorTypeCommunication:  model.NewErrCommunication,
	model.ErrorTypeAuthorization:  model.NewErrAuthorization,
	model.ErrorTypeConfiguration:  model.NewErrConfiguration,
	model.ErrorTypeExpression:     model.NewErrExpression,
	model.ErrorTypeRuntime:        model.NewErrRuntime,
	model.ErrorTypeTimeout:        model.NewErrTimeout,
}

func (r *RaiseTaskRunner) Run(input interface{}) (output interface{}, err error) {
	output = input
	// TODO: make this an external func so we can call it after getting the reference? Or we can get the reference from the workflow definition
	var detailResult interface{}
	detailResult, err = traverseAndEvaluate(r.Task.Raise.Error.Definition.Detail.AsObjectOrRuntimeExpr(), input, r.TaskName)
	if err != nil {
		return nil, err
	}

	var titleResult interface{}
	titleResult, err = traverseAndEvaluate(r.Task.Raise.Error.Definition.Title.AsObjectOrRuntimeExpr(), input, r.TaskName)
	if err != nil {
		return nil, err
	}

	instance := &model.JsonPointerOrRuntimeExpression{Value: r.TaskName}

	var raiseErr *model.Error
	if raiseErrF, ok := raiseErrFuncMapping[r.Task.Raise.Error.Definition.Type.String()]; ok {
		raiseErr = raiseErrF(fmt.Errorf("%v", detailResult), instance.String())
	} else {
		raiseErr = r.Task.Raise.Error.Definition
		raiseErr.Detail = model.NewStringOrRuntimeExpr(fmt.Sprintf("%v", detailResult))
		raiseErr.Instance = instance
	}

	raiseErr.Title = model.NewStringOrRuntimeExpr(fmt.Sprintf("%v", titleResult))
	err = raiseErr

	return output, err
}

func (r *RaiseTaskRunner) GetTaskName() string {
	return r.TaskName
}
