package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/lpabon/godbc"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/oauth2"
)

type encrypted struct {
	Ciphertext []byte
	Nonce      []byte
}

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

func encryptToken(tok *oauth2.Token) []byte {
	// turn token to bytes, then feed to encrypt
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(tok); err != nil {
		panic(err.Error())
	}
	out := buf.Bytes()
	godbc.Ensure(len(out) > 0, "didn't create out bytes")
	return out
}

func decryptToken(b []byte) *oauth2.Token {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	tok := new(oauth2.Token)
	if err := dec.Decode(tok); err != nil {
		panic(err.Error())
	}
	return tok
}

func encrypt(b []byte) []byte {
	godbc.Require(len(b) > 0, "input length must be > 0")

	block, err := aes.NewCipher(getKey())
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, 12)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	crypt := encrypted{gcm.Seal(nil, nonce, b, nil), nonce}
	if err := enc.Encode(crypt); err != nil {
		panic(err.Error())
	}
	out := buf.Bytes()
	godbc.Ensure(len(out) > 0, "didn't create out bytes")
	return out
}

func decrypt(buf []byte) []byte {
	godbc.Require(len(buf) > 0, "buf length must be > 0")

	var crypt encrypted
	dec := gob.NewDecoder(bytes.NewBuffer(buf))
	if err := dec.Decode(&crypt); err != nil {
		panic(err.Error())
	}

	block, err := aes.NewCipher(getKey())
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	plaintext, err := gcm.Open(nil, crypt.Nonce, crypt.Ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func saveToken(tok *oauth2.Token, path string) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		panic(err.Error())
	}

	if _, err := f.Write(encryptToken(tok)); err != nil {
		panic(err.Error())
	}
}

func loadToken(path string) *oauth2.Token {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		// fuck it, just return nil
		return nil
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, f); err != nil {
		panic(err.Error())
	}

	tok := decryptToken(buf.Bytes())
	if tok.Expiry.Before(time.Now()) {
		return nil
	}
	return tok
}
