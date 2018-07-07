package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeAndDecode(t *testing.T) {
	ct, nonce := encode([]byte("hi"))
	assert.Equal(t, []byte("hi"), decode(ct, nonce))
}
