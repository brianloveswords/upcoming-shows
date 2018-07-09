package main

import (
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/lpabon/godbc"
)

func mustGet(v string, key string) string {
	if envVal := os.Getenv(key); envVal != "" {
		glog.Extreme("overriding from %s", key)
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
		glog.Debug("couldn't open intmap data from file '%s'", name)
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
		fmt.Printf("go here\n: %s\n", url)
		return
	}
}
