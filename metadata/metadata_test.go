package metadata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetadata_Get(t *testing.T) {
	key := "key1"
	value := "value1"

	md := Metadata(map[string]string{
		key: value,
	})

	assert.Equal(t, value, md.Get("KEY1"))
}

func TestMetadata_Set(t *testing.T) {
	key := "key1"
	value := "value1"

	md := Metadata(make(map[string]string))

	md.Set(key, value)

	assert.Equal(t, value, md[key])
}

func TestMetadata_Map(t *testing.T) {
	key := "key1"
	key2 := "key2"
	value := "value1"

	md := Metadata(map[string]string{
		key: value,
	})

	md.Map(key, key2)

	assert.Equal(t, value, md[key2])
}
