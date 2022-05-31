package encoder

import (
	"bytes"
	"fmt"
	"strings"
)

type BinaryEncoder struct {
}

func (e *BinaryEncoder) String(src []byte) (string, error) {

	var builder strings.Builder
	for _, b := range src {
		builder.WriteString(byteToHex(b))
	}

	return builder.String(), nil
}

func (e *BinaryEncoder) Bytes(src string) ([]byte, error) {

	var buf bytes.Buffer

	// 3文字で1つのヘキサ文字(x00)になっている
	for i := 0; i < len(src); i += 3 {

		end := i + 3
		if end > len(src) {
			end = len(src)
		}

		b, err := hexToByte(src[i:end])
		if err != nil {
			return nil, err
		}

		buf.WriteByte(b)
	}

	return buf.Bytes(), nil
}

const hextable = "0123456789ABCDEF"

func byteToHex(b byte) string {

	return string([]byte{'x', hextable[b/16], hextable[b%16]})
}

func hexToByte(h string) (byte, error) {

	if len(h) != 3 {
		return 0x00, fmt.Errorf("illegal hex string \"%s\"", h)
	}

	if h[0] != 'x' {
		return 0x00, fmt.Errorf("illegal hex string \"%s\"", h)
	}

	first := strings.IndexByte(hextable, h[1])
	if first == -1 {
		return 0x00, fmt.Errorf("illegal hex string \"%s\"", h)
	}

	second := strings.IndexByte(hextable, h[2])
	if second == -1 {
		return 0x00, fmt.Errorf("illegal hex string \"%s\"", h)
	}

	return byte(first*16 + second), nil
}
