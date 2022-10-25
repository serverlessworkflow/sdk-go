package parser

import (
	"embed"
	"fmt"
	"sort"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/workflows/*
var contentFiles embed.FS

func Test_parseSchema(t *testing.T) {
	path := "testdata/workflows"
	sch, err := jsonschema.Compile("https://serverlessworkflow.io/schemas/0.8/workflow.json")
	assert.NoError(t, err)

	err = validateSchema(contentFiles, path, sch)
	if merr, ok := err.(*multierror.Error); ok {
		fmt.Printf("There are %d errors to solve", merr.Len())
		var errs []string
		for i := range merr.Errors {
			errs = append(errs, merr.Errors[i].Error())
		}
		sort.Strings(errs)

		for i := range errs {
			fmt.Println(errs[i])
		}
	}
	assert.Nil(t, err)
}
