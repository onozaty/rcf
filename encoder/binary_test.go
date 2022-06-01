package encoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestByteToHex(t *testing.T) {

	assert.Equal(t, "x00", byteToHex(0x00))
	assert.Equal(t, "x01", byteToHex(0x01))
	assert.Equal(t, "x10", byteToHex(0x10))
	assert.Equal(t, "x99", byteToHex(0x99))
	assert.Equal(t, "x0A", byteToHex(0x0A))
	assert.Equal(t, "xAA", byteToHex(0xAA))
	assert.Equal(t, "xBC", byteToHex(0xBC))
	assert.Equal(t, "xDE", byteToHex(0xDE))
	assert.Equal(t, "xF0", byteToHex(0xF0))
	assert.Equal(t, "xFF", byteToHex(0xFF))
}

func TestHexToByte(t *testing.T) {

	{
		b, err := hexToByte("x00")
		assert.NoError(t, err)
		assert.Equal(t, byte(0x00), b)
	}
	{
		b, err := hexToByte("x01")
		assert.NoError(t, err)
		assert.Equal(t, byte(0x01), b)
	}
	{
		b, err := hexToByte("x10")
		assert.NoError(t, err)
		assert.Equal(t, byte(0x10), b)
	}
	{
		b, err := hexToByte("x99")
		assert.NoError(t, err)
		assert.Equal(t, byte(0x99), b)
	}
	{
		b, err := hexToByte("x0A")
		assert.NoError(t, err)
		assert.Equal(t, byte(0x0A), b)
	}
	{
		b, err := hexToByte("xAA")
		assert.NoError(t, err)
		assert.Equal(t, byte(0xAA), b)
	}
	{
		b, err := hexToByte("xBC")
		assert.NoError(t, err)
		assert.Equal(t, byte(0xBC), b)
	}
	{
		b, err := hexToByte("xDE")
		assert.NoError(t, err)
		assert.Equal(t, byte(0xDE), b)
	}
	{
		b, err := hexToByte("xF0")
		assert.NoError(t, err)
		assert.Equal(t, byte(0xF0), b)
	}
	{
		b, err := hexToByte("xFF")
		assert.NoError(t, err)
		assert.Equal(t, byte(0xFF), b)
	}
}

func TestHexToByte_Invalid(t *testing.T) {

	{
		_, err := hexToByte("FF")
		assert.Error(t, err)
		assert.Equal(t, `illegal hex string "FF"`, err.Error())
	}

	{
		_, err := hexToByte("XFF")
		assert.Error(t, err)
		assert.Equal(t, `illegal hex string "XFF"`, err.Error())
	}

	{
		_, err := hexToByte("xff")
		assert.Error(t, err)
		assert.Equal(t, `illegal hex string "xff"`, err.Error())
	}

	{
		_, err := hexToByte("x0f")
		assert.Error(t, err)
		assert.Equal(t, `illegal hex string "x0f"`, err.Error())
	}
}

func TestNewEncoder_Binary(t *testing.T) {

	// ARRANGE
	str := "x00x01x70x71x80x81xF0xFF"
	bytes := []byte{'\x00', '\x01', '\x70', '\x71', '\x80', '\x81', '\xF0', '\xFF'}

	// ACT / ASSERT
	encoder, err := NewEncoder("binary")
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

func TestNewEncoder_Binary_Invalid(t *testing.T) {

	// ARRANGE
	str := "x00x01x70x71x80x81xF0xF"

	// ACT / ASSERT
	encoder, err := NewEncoder("binary")
	require.NoError(t, err)

	_, err = encoder.Bytes(str)
	require.Error(t, err)
	assert.Equal(t, `illegal hex string "xF"`, err.Error())
}
