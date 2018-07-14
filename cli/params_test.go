package cli

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParamBase(t *testing.T) {
	b := Base{
		Name:  "test",
		Short: 't',
	}

	cases := []struct {
		input       string
		expectfound bool
		expectvalue string
	}{
		{"--test=10", true, "10"},
		{"--test", true, ""},
		{"-t=10", true, "10"},
		{"-t", true, ""},

		{"--testd=1", false, ""},
		{"-test=1", false, ""},
		{"-test1", false, ""},
		{"-n=1", false, ""},
		{"-t1", false, ""},
	}

	for _, tt := range cases {
		t.Run(tt.input, func(t *testing.T) {
			value, found := b.Extract(tt.input)
			assert.Equal(t, tt.expectfound, found)
			assert.Equal(t, tt.expectvalue, value)
		})
	}
}

func TestParamParser(t *testing.T) {
	var (
		// paramLength  *int
		// paramBool    *bool
		// paramTrackID *spotify.ID
		paramArtist *string
	)

	pl := ParamList{
		&StringParam{
			Base: Base{
				Name:  "artist",
				Help:  "artist name to make mix from",
				Short: 'a',
			},
			Var:     paramArtist,
			Default: "%currentartist%",
		},
	}

	pl.Parse([]string{"hi"})
	fmt.Println(paramArtist)

	pl.Parse([]string{"--artist=hi"})
	fmt.Println(paramArtist)

	pl.Parse([]string{"-aRad"})
	fmt.Println(paramArtist)

	pl.Parse([]string{"-a=Rad"})
	fmt.Println(paramArtist)
}
