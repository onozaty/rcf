package encoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEncoder_UTF8(t *testing.T) {

	// ARRANGE
	str := "あいうえお"
	bytes := []byte{'\xE3', '\x81', '\x82', '\xE3', '\x81', '\x84', '\xE3', '\x81', '\x86', '\xE3', '\x81', '\x88', '\xE3', '\x81', '\x8A'}

	// ACT / ASSERT
	encoder, err := NewEncoder("utf-8")
	require.NoError(t, err)

	{
		result, err := encoder.String(bytes)
		require.NoError(t, err)
		assert.Equal(t, str, result)
	}

	{
		result, err := encoder.Bytes(str)
		require.NoError(t, err)
		assert.Equal(t, bytes, result)
	}
}

func TestNewEncoder_SJIS(t *testing.T) {

	// ARRANGE
	str := "あいうえお"
	bytes := []byte{'\x82', '\xA0', '\x82', '\xA2', '\x82', '\xA4', '\x82', '\xA6', '\x82', '\xA8'}

	// ACT / ASSERT
	encoder, err := NewEncoder("sjis")
	require.NoError(t, err)

	{
		result, err := encoder.String(bytes)
		require.NoError(t, err)
		assert.Equal(t, str, result)
	}

	{
		result, err := encoder.Bytes(str)
		require.NoError(t, err)
		assert.Equal(t, bytes, result)
	}
}

func TestNewEncoder_Invalid(t *testing.T) {

	// ACT / ASSERT
	_, err := NewEncoder("xxxx")
	require.Error(t, err)
	assert.Equal(t, "htmlindex: invalid encoding name", err.Error())
}
