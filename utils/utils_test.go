package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase58(t *testing.T) {
	info := "nihao"
	encodeInfo := Base58Encode([]byte(info))
	decodeInfo := Base58Decode(encodeInfo)
	fmt.Println(encodeInfo, decodeInfo, string(decodeInfo))
	assert.Equal(t, info, string(decodeInfo))
}
