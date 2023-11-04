package util

import (
	"encoding/json"
	"reflect"

	"github.com/jmoiron/sqlx/types"
)

// CompareJSONNullText compares two types.NullJSONText objects for equality, ignoring key order.
func CompareNullJSONText(json1, json2 types.NullJSONText) (bool, error) {
	if json1.Valid != json2.Valid {
		return false, nil
	}
	if !json1.Valid {
		// both are null so they match
		return true, nil
	}

	var obj1, obj2 any
	// Unmarshal the json into an any.
	err := json.Unmarshal(json1.JSONText, &obj1)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(json2.JSONText, &obj2)
	if err != nil {
		return false, err
	}

	// Compare the resulting any objects.
	return reflect.DeepEqual(obj1, obj2), nil
}
