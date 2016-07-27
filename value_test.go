package failsafe

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValues(t *testing.T) {
	v := newValue()

	v.Put("hello", 1)
	assert.Equal(t, v.Len(), 1)

	v.Put("hello2", "earth")
	assert.Equal(t, v.Len(), 2)

	v.Put("hello2", "world")
	assert.Equal(t, v.Len(), 2)

	val, ok := v.Get("hello2")
	assert.True(t, ok)
	assert.Equal(t, val.String(), "world")

	val, ok = v.Get("hello")
	assert.True(t, ok)
	assert.Equal(t, val.Int(), 1)

	intSlice := []int{1, 2, 3, 4, 5}
	v.Put("slice", intSlice)
	assert.Equal(t, v.Len(), 3)
	val, ok = v.Get("slice")
	assert.True(t, ok)
	assert.Equal(t, len(val.Slice()), len(intSlice))
	for i, v := range intSlice {
		assert.Equal(t, val.Slice()[i].Int(), v)
	}

	v.Remove("hello2")
	assert.Equal(t, v.Len(), 2)
	val, ok = v.Get("hello2")
	assert.False(t, ok)
	assert.Equal(t, val, (*abstractValue)(nil))

	v.Clear()
	assert.Equal(t, v.Len(), 0)
}
