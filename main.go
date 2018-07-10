package main

import (
	"os"

	"github.com/brianloveswords/spotify/auth"
	"github.com/brianloveswords/spotify/cli"
	"github.com/brianloveswords/spotify/logger"
	"github.com/brianloveswords/spotify/mixtape"
	"github.com/brianloveswords/spotify/show"
	"github.com/brianloveswords/spotify/util"
	"github.com/fatih/color"
)

// these will be set at build time from env SPOTIFY_ID and SPOTIFY_SECRET
var clientID string
var clientSecret string

var glog = logger.DefaultLogger

func cliRouter(args []string) {
	mainCmd := cli.Command{
		Name: "spotify",
		Help: "control spotify from the commandline",
		Commands: cli.Subcommands{
			&cli.Command{
				Name: "play",
				Help: "play the current song",
				Func: mainPlay,
			},
			&cli.Command{
				Name: "pause",
				Help: "pause the current song",
				Func: mainPause,
			},
			&cli.Command{
				Name:  "skip",
				Help:  "skip the current song",
				Alias: []string{"next"},
				Func:  mainNext,
			},
			&cli.Command{
				Name:  "prev",
				Help:  "go back to the last song",
				Alias: []string{"back"},
				Func:  mainPrev,
			},
			&cli.Command{
				Name: "fav",
				Help: "add the current song to your library",
				Func: mainFav,
			},
			&show.CLI,
			&mixtape.CLI,
		},
	}
	err := mainCmd.Run(args)
	if err != nil {
		glog.Fatal("%s", err)
	}
}
func mainFav() {
	client := auth.SetupClient()
	playing, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		glog.Fatal("could not get currently playing: %s", err)
	}
	track := playing.Item
	if err := client.AddTracksToLibrary(track.ID); err != nil {
		glog.Fatal("could add track to library: %s", err)
	}
	name := util.SongAttributionFromTrack(track)
	glog.Log("adding to library: %s", color.CyanString(name))
}

func mainPlay() {
	client := auth.SetupClient()
	if glog.IsLevelNormal() {
		util.LogCurrentTrack(client, "playing")
	}
	if err := client.Play(); err != nil {
		glog.Fatal("couldn't start playback: %s", err)
	}
}
func mainPause() {
	client := auth.SetupClient()
	if glog.IsLevelNormal() {
		util.LogCurrentTrack(client, "pausing")
	}
	if err := client.Pause(); err != nil {
		glog.Fatal("couldn't pause playback: %s", err)
	}
}

func mainNext() {
	client := auth.SetupClient()
	util.LogCurrentTrack(client, "skipping")
	if err := client.Next(); err != nil {
		glog.Fatal("couldn't skip track: ", err)
	}
}
func mainPrev() {
	client := auth.SetupClient()
	if err := client.Previous(); err != nil {
		glog.Fatal("couldn't go back: ", err)
	}
	util.LogCurrentTrack(client, "going back to")
}

func main() {
	cliRouter(os.Args[1:])
}
