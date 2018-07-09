package main

import (
	"fmt"
	"testing"
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

	mixtapeCmds = Command{
		Name: "mixtape",
		Help: "show details about the current track",
		Commands: Subcommands{
			&Command{
				Name: "current-artist",
				Help: "create a mixtape playlist for the current artist",
			},
		},
	}

	fmt.Println(mainCmd.String())
	mainCmd.Run([]string{"show", "track-uri"})
}
