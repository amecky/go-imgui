package imgui

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestFindMaxLen(t *testing.T) {
	l := findMaxLen([]string{"1", "12", "123"})
	assert.Equal(t, 3, l)

	l = findMaxLen([]string{"1"})
	assert.Equal(t, 1, l)

	l = findMaxLen([]string{"1234", "12", "123"})
	assert.Equal(t, 4, l)

}
