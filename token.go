package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/scrypt"
)

func token() {
	fmt.Println("lol")
}

func getKey() []byte {
	id, ok := os.LookupEnv("SPOTIFY_ID")
	if !ok {
		panic("SPOTIFY_ID not set")
	}
	secret, ok := os.LookupEnv("SPOTIFY_SECRET")
	if !ok {
		panic("SPOTIFY_SECRET not set")
	}
	dk, err := scrypt.Key([]byte(secret), []byte(id), 32768, 8, 1, 32)
	if err != nil {
		panic(err)
	}
	return dk
}

func encode(b []byte) (ciphertext, nonce []byte) {
	block, err := aes.NewCipher(getKey())
	if err != nil {
		panic(err.Error())
	}
	nonce = make([]byte, 12)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	return gcm.Seal(nil, nonce, b, nil), nonce
}

func decode(ciphertext, nonce []byte) []byte {
	block, err := aes.NewCipher(getKey())
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}
