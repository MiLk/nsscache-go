package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullInt32_String(t *testing.T) {
	n := nullInt32{}
	assert.Empty(t, n.String())

	v := Int32(200)
	assert.Equal(t, "200", v.String())
}

func TestNullUInt32_String(t *testing.T) {
	n := nullUInt32{}
	assert.Empty(t, n.String())

	v := UInt32(200)
	assert.Equal(t, "200", v.String())
}
