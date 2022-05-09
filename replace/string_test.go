package replace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringReplacer(t *testing.T) {

	replacer := NewStringReplacer("abc", "xyz")

	{
		result := replacer.Replace("abc")
		assert.Equal(t, "xyz", result)
	}
	{
		result := replacer.Replace(" abcabc\nabcABC\n")
		assert.Equal(t, " xyzxyz\nxyzABC\n", result)
	}
	{
		result := replacer.Replace("")
		assert.Equal(t, "", result)
	}
	{
		result := replacer.Replace("aaaa")
		assert.Equal(t, "aaaa", result)
	}
}
