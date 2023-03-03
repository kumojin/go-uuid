package uuid

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_StringArray_Scan(t *testing.T) {
	var (
		uuid1, _  = FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		uuid2, _  = FromString("3964fdbc-dc5a-4ec6-8a6b-eecfa5a8e9ae")
		uuidArray = UUIDArray{uuid1, uuid2}
	)

	json := `["6ba7b810-9dad-11d1-80b4-00c04fd430c8","3964fdbc-dc5a-4ec6-8a6b-eecfa5a8e9ae"]`
	newArrays := UUIDArray{}
	err := newArrays.Scan(json)
	assert.Equal(t, uuidArray, newArrays)
	assert.NoError(t, err)

	newArrays = UUIDArray{}
	err = newArrays.Scan([]byte(json))
	assert.Equal(t, uuidArray, newArrays)
	assert.NoError(t, err)

	err = newArrays.Scan(nil)
	assert.Equal(t, newArrays, UUIDArray{})
	assert.NoError(t, err)

	err = newArrays.Scan(42)
	assert.Error(t, err)

	err = newArrays.Scan("{I am not a valid json")
	assert.Error(t, err)

	err = newArrays.Scan(nil)
	assert.NoError(t, err)
}

func Test_StringArray_Value(t *testing.T) {
	var (
		uuid1, _  = FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
		uuid2, _  = FromString("3964fdbc-dc5a-4ec6-8a6b-eecfa5a8e9ae")
		uuidArray = UUIDArray{uuid1, uuid2}
	)

	expected := `["6ba7b810-9dad-11d1-80b4-00c04fd430c8","3964fdbc-dc5a-4ec6-8a6b-eecfa5a8e9ae"]`
	val, _ := uuidArray.Value()
	assert.Equal(t, expected, fmt.Sprintf("%s", val))
}
