package builder

import (
	"encoding/json"

	"sigs.k8s.io/yaml"

	"github.com/serverlessworkflow/sdk-go/v2/model"
)

func New() *model.WorkflowBuilder {
	return model.NewWorkflowBuilder()
}

func AsObject(builder *model.WorkflowBuilder) *model.Workflow {
	workflow := builder.Build()
	return &workflow
}

func AsJson(builder *model.WorkflowBuilder) ([]byte, error) {
	workflow := builder.Build()
	return json.Marshal(workflow)
}

func AsYaml(builder *model.WorkflowBuilder) ([]byte, error) {
	data, err := AsJson(builder)
	if err != nil {
		return nil, err
	}

	return yaml.JSONToYAML(data)
}
