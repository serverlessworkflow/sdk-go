package parser

import (
	"embed"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed testdata/workflows/*
var contentFiles embed.FS

func Test_parseSchema(t *testing.T) {
	path := "testdata/workflows"
	sch, err := jsonschema.Compile("testdata/schema/workflows.json")
	assert.NoError(t, err)

	err = validateSchema(contentFiles, path, sch)
	assert.NoError(t, err)
}
