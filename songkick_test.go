package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSongkickParser(t *testing.T) {
	f, _ := os.Open("sk2.html")
	defer f.Close()
	fmt.Println(getIDFromSongkickPage("Gleemer"))
}

func TestSongkickIDFromHref(t *testing.T) {
	assert.Equal(t, 7180534, idFromHref("/artists/7180534-gleemer"))
}
