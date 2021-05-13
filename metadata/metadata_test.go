package metadata

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadata_Get(t *testing.T) {
	key := "key1"
	value := "value1"

	md := Metadata(map[string]string{
		key: value,
	})

	assert.Equal(t, value, md.GetRaw("KEY1"))
}

func TestMetadata_Set(t *testing.T) {
	key := "key1"
	value := "value1"

	md := Metadata(make(map[string]string))

	md.set(key, value)

	assert.Equal(t, value, md[key])
}

func TestMetadata_ConvertToType(t *testing.T) {
	tests := []struct {
		name string
		m    Metadata
		want UploadMetadata
	}{
		// TODO: Add test cases.
		{
			"Should not fail on empty values",
			Metadata{},
			UploadMetadata{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.GetUploadMetadata(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUploadMetadata() = \n%+v, \nwant \n%+v", got, tt.want)
			}
		})
	}
}
