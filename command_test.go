package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/zmb3/spotify"
)

func TestCommandSmokeTest(t *testing.T) {
	var (
		mainCmd     Command
		mixtapeCmds Command
		showCmds    Command
	)

	mainCmd = Command{
		Name: "spotify",
		Help: "control spotify from the commandline",
		Commands: Subcommands{
			&Command{
				Name: "play",
				Help: "play the current song",
			},
			&Command{
				Name: "pause",
				Help: "pause the current song",
			},
			&Command{
				Name:  "skip",
				Help:  "skip the current song",
				Alias: []string{"next"},
			},
			&Command{
				Name:  "prev",
				Help:  "go back to the last song",
				Alias: []string{"back"},
			},
			&showCmds,
			&mixtapeCmds,
		},
	}

	showCmds = Command{
		Name: "show",
		Help: "show details about the current track",
		Commands: Subcommands{
			&Command{
				Name: "artist",
				Help: "show the artist of the current track",
			},
			&Command{
				Name: "artist-id",
				Help: "show the spotify artist ID of the current track",
			},
			&Command{
				Name: "artist-uri",
				Help: "show the spotify URI for the artist of the current track",
			},
			&Command{
				Name: "track",
				Help: "show the title for the current track",
			},
			&Command{
				Name: "track-id",
				Help: "show the track ID of the current track",
			},
			&Command{
				Name: "track-uri",
				Help: "show the spotify URI for the current track",
				Fn: func() {
					fmt.Println("gonna show that track-uri")
				},
			},
		},
	}
	var (
		paramMixtapeArtist   string
		paramMixtapeTrackID  spotify.ID
		paramMixtapeLength   int
		defaultMixtapeLength = 10
	)
	mixtapeCmds = Command{
		Name: "mixtape",
		Help: "create a mixtape",
		Fn: func() {
			fmt.Printf("artist is %s\n", paramMixtapeArtist)
			fmt.Printf("track is %s\n", paramMixtapeTrackID)
			fmt.Printf("length is %d\n", paramMixtapeLength)
		},
		Params: []Param{
			Param{
				Name:  "artist",
				Alias: []string{"a"},
				Help:  "the artist to base the mixtape on. If passed and not set, defaults to currently playing artist. Mutually exclusive with `track`",
				ParseFn: func(val string) error {
					if paramMixtapeTrackID != "" {
						return fmt.Errorf("must pass track or artist, but not both")
					}
					paramMixtapeArtist = val
					return nil
				},
			},
			Param{
				Name:  "track",
				Alias: []string{"t"},
				Help:  "the track ID to base the mixtape on. If passed and not set, defaults to currently playing artist. Mutually exclusive with `artist`",
				ParseFn: func(val string) error {
					if paramMixtapeArtist != "" {
						return fmt.Errorf("must pass track or artist, but not both")
					}
					paramMixtapeTrackID = spotify.ID(val)
					return nil
				},
			},
			Param{
				Name:     "length",
				Alias:    []string{"n"},
				Help:     fmt.Sprintf("number of tracks to. defaults to %d", defaultMixtapeLength),
				Implicit: true,
				ParseFn: func(val string) (err error) {
					if val == "" {
						paramMixtapeLength = defaultMixtapeLength
						return nil
					}
					paramMixtapeLength, err = strconv.Atoi(val)
					return err
				},
			},
		},
	}

	var err error
	fmt.Println(mainCmd.String())
	mainCmd.Run([]string{"show", "track-uri"})
	err = mainCmd.Run([]string{"mixtape", "artist=chavez"})
	if err != nil {
		fmt.Printf("err %s\n", err)
	}
	err = mainCmd.Run([]string{"mixtape", "n=20", "artist=chavez"})
	if err != nil {
		fmt.Printf("err %s\n", err)
	}
	err = mainCmd.Run([]string{"mixtape", "a=chavez"})
	if err != nil {
		fmt.Printf("err %s\n", err)
	}
	err = mainCmd.Run([]string{"mixtape", "a=chavez", "t=a2f"})
	if err != nil {
		fmt.Printf("err %s\n", err)
	}
}
