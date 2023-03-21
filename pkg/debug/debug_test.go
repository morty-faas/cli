package debug

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_JSON(t *testing.T) {
	expected := `{
	"foo": "bar"
}`
	assert.Equal(t, expected, JSON(map[string]interface{}{
		"foo": "bar",
	}))
}
