package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestBasicEncryption(t *testing.T) {
	assert.Equal(t, []byte("hi"), decrypt(encrypt([]byte("hi"))))
}

func TestEncryptDecryptToken(t *testing.T) {
	tok := &oauth2.Token{
		AccessToken:  "accesstoken",
		TokenType:    "tokentype",
		RefreshToken: "refreshtoken",
		Expiry:       time.Now(),
	}

	tok2 := decryptToken(encryptToken(tok))

	assert.Equal(t, tok.AccessToken, tok2.AccessToken)
	assert.Equal(t, tok.TokenType, tok2.TokenType)
	assert.Equal(t, tok.RefreshToken, tok2.RefreshToken)
}
