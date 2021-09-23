package floatstr

import (
	"encoding/json"
	"k8s.io/apimachinery/pkg/util/yaml"
	"reflect"
	"testing"
)

func TestFromFloat(t *testing.T) {
	i := FromFloat(93.93)
	if i.Type != Float || i.FloatVal != 93.93 {
		t.Errorf("Expected FloatVal=93.93, got %+v", i)
	}
}

func TestFromString(t *testing.T) {
	i := FromString("76.76")
	if i.Type != String || i.StrVal != "76.76" {
		t.Errorf("Expected StrVal=\"76.76\", got %+v", i)
	}
}

type FloatOrStringHolder struct {
	FOrS Float32OrString `json:"val"`
}

func TestIntOrStringUnmarshalJSON(t *testing.T) {
	cases := []struct {
		input  string
		result Float32OrString
	}{
		{"{\"val\": 123.123}", FromFloat(123.123)},
		{"{\"val\": \"123.123\"}", FromString("123.123")},
	}

	for _, c := range cases {
		var result FloatOrStringHolder
		if err := json.Unmarshal([]byte(c.input), &result); err != nil {
			t.Errorf("Failed to unmarshal input '%v': %v", c.input, err)
		}
		if result.FOrS != c.result {
			t.Errorf("Failed to unmarshal input '%v': expected %+v, got %+v", c.input, c.result, result)
		}
	}
}

func TestIntOrStringMarshalJSON(t *testing.T) {
	cases := []struct {
		input  Float32OrString
		result string
	}{
		{FromFloat(123.123), "{\"val\":123.123}"},
		{FromString("123.123"), "{\"val\":\"123.123\"}"},
	}

	for _, c := range cases {
		input := FloatOrStringHolder{c.input}
		result, err := json.Marshal(&input)
		if err != nil {
			t.Errorf("Failed to marshal input '%v': %v", input, err)
		}
		if string(result) != c.result {
			t.Errorf("Failed to marshal input '%v': expected: %+v, got %q", input, c.result, string(result))
		}
	}
}

func TestIntOrStringMarshalJSONUnmarshalYAML(t *testing.T) {
	cases := []struct {
		input Float32OrString
	}{
		{FromFloat(123.123)},
		{FromString("123.123")},
	}

	for _, c := range cases {
		input := FloatOrStringHolder{c.input}
		jsonMarshalled, err := json.Marshal(&input)
		if err != nil {
			t.Errorf("1: Failed to marshal input: '%v': %v", input, err)
		}

		var result FloatOrStringHolder
		err = yaml.Unmarshal(jsonMarshalled, &result)
		if err != nil {
			t.Errorf("2: Failed to unmarshal '%+v': %v", string(jsonMarshalled), err)
		}

		if !reflect.DeepEqual(input, result) {
			t.Errorf("3: Failed to marshal input '%+v': got %+v", input, result)
		}
	}
}
