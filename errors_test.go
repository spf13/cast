package cast

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLazyError_Error(t *testing.T) {
	err := newError("foo")

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "foo")
}
