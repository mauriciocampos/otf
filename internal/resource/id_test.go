package resource

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	id := NewTfeID("foo")
	got := fmt.Sprintf("%v", id)
	assert.Regexp(t, `foo-.*`, got)
}
