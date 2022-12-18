package utils_test

import (
	"github.com/go-playground/assert/v2"
	"github.com/miniyus/go-fiber/internal/utils"
	"testing"
)

func TestB64UrlEncode(t *testing.T) {
	testString := "hello, world"

	encodeString := utils.B64UrlEncode(testString)
	decodeString, err := utils.B64UrlDecode(encodeString)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, decodeString, testString)
}
