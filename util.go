package main

import (
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

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

func mustCreate(filename string) *os.File {
	f, err := os.Create(filename)
	if err != nil {
		panic(err.Error())
	}
	return f
}

func loadIntMap(name string) (intmap map[string]int) {
	f, err := os.Open(name)
	if err != nil {
		debugprint("couldn't open intmap data from file '%s'\n", name)
		return make(map[string]int)
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&intmap); err != nil {
		panic(err.Error())
	}
	return intmap
}

func saveIntMap(name string, intmap map[string]int) {
	f := mustCreate(name)
	defer f.Close()

	enc := gob.NewEncoder(f)
	if err := enc.Encode(intmap); err != nil {
		panic(err)
	}
}

func openURL(url string) {
	cmd := exec.Command("open", url)
	if err := cmd.Run(); err != nil {
		// fall back to just printing it
		fmt.Printf("go here and authenticate\n: %s\n", url)
		return
	}
}
