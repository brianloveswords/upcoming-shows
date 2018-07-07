package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/lpabon/godbc"
)

var dbglog = log.New(os.Stderr, "debug: ", 0)

func debugprint(format string, v ...interface{}) {
	if _, ok := os.LookupEnv("DEBUG"); ok {
		dbglog.Printf(format, v...)
	}
}
func mustGet(v string, key string) string {
	if envVal := os.Getenv(key); envVal != "" {
		debugprint("overriding from %s", key)
		v = envVal
	}
	godbc.Ensure(v != "", fmt.Sprintf("couldn't ensure value with key %s", key))
	return v
}

func randomState() string {
	b := make([]byte, 24)
	io.ReadFull(rand.Reader, b)
	return hex.EncodeToString(b)
}
