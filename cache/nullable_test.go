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

func TestNullInt32_MarshalJSON(t *testing.T) {
	n := nullInt32{}
	b, err := n.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, []byte(`""`), b)

	v := Int32(200)
	b, err = v.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, []byte(`"200"`), b)

}

func TestNullInt32_UnmarshalJSON(t *testing.T) {
	n := nullInt32{}
	assert.Nil(t, n.UnmarshalJSON([]byte(`""`)))
	assert.False(t, n.valid)

	assert.Nil(t, n.UnmarshalJSON([]byte(`"200"`)))
	assert.True(t, n.valid)
	assert.EqualValues(t, 200, n.val)

	assert.NotNil(t, n.UnmarshalJSON([]byte(`['([)]'/`)))
	assert.NotNil(t, n.UnmarshalJSON([]byte(`"foo"`)))
}

func TestNullUInt32_String(t *testing.T) {
	n := nullUInt32{}
	assert.Empty(t, n.String())

	v := UInt32(200)
	assert.Equal(t, "200", v.String())
}

func TestNullUInt32_MarshalJSON(t *testing.T) {
	n := nullUInt32{}
	b, err := n.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, []byte(`""`), b)

	v := UInt32(200)
	b, err = v.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, []byte(`"200"`), b)

}

func TestNullUInt32_UnmarshalJSON(t *testing.T) {
	n := nullUInt32{}
	assert.Nil(t, n.UnmarshalJSON([]byte(`""`)))
	assert.False(t, n.valid)

	assert.Nil(t, n.UnmarshalJSON([]byte(`"200"`)))
	assert.True(t, n.valid)
	assert.EqualValues(t, 200, n.val)

	assert.NotNil(t, n.UnmarshalJSON([]byte(`['([)]'/`)))
	assert.NotNil(t, n.UnmarshalJSON([]byte(`"foo"`)))
}
