package util_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CompareJSON compares two json.RawMessage objects for equality, ignoring key order.
func CompareJSON(json1, json2 json.RawMessage) (bool, error) {
	var obj1, obj2 any

	// Unmarshal the json into an any.
	err := json.Unmarshal(json1, &obj1)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(json2, &obj2)
	if err != nil {
		return false, err
	}

	// Compare the resulting any objects.
	return reflect.DeepEqual(obj1, obj2), nil
}

func TestCompareJSON(t *testing.T) {
	tests := []struct {
		name  string
		json1 string
		json2 string
		want  bool
	}{
		{
			name:  "empty",
			json1: `{}`,
			json2: `{}`,
			want:  true,
		},
		{
			name:  "simple",
			json1: `{"a": "b"}`,
			json2: `{"a": "b"}`,
			want:  true,
		},
		{
			name:  "simple different",
			json1: `{"a": "b"}`,
			json2: `{"a": "c"}`,
			want:  false,
		},
		{
			name:  "nested",
			json1: `{"a": {"b": "c"}}`,
			json2: `{"a": {"b": "c"}}`,
			want:  true,
		},
		{
			name:  "nested different",
			json1: `{"a": {"b": "c"}}`,
			json2: `{"a": {"b": "d"}}`,
			want:  false,
		},
		{
			name:  "nested different key order",
			json1: `{"a": {"b": "c", "d": "e"}}`,
			json2: `{"a": {"d": "e", "b": "c"}}`,
			want:  true,
		},
		{
			name:  "nested different key order",
			json1: `{"a": {"b": "c", "d": "e"}}`,
			json2: `{"a": {"d": "e", "b": "f"}}`,
			want:  false,
		},
		{
			name:  "array",
			json1: `{"a": ["b", "c"]}`,
			json2: `{"a": ["b", "c"]}`,
			want:  true,
		},
		{
			name:  "array different",
			json1: `{"a": ["b", "c"]}`,
			json2: `{"a": ["b", "d"]}`,
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			json1 := json.RawMessage(tt.json1)
			json2 := json.RawMessage(tt.json2)
			got, err := CompareJSON(json1, json2)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
