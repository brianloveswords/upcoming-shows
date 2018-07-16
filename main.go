package main

import (
	"os"

	"github.com/brianloveswords/spotify/auth"
	"github.com/brianloveswords/spotify/logger"
	"github.com/brianloveswords/spotify/util"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var glog = logger.DefaultLogger

func mainFav(c *cli.Context) error {
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
	return nil
}

func mainPlay(c *cli.Context) error {
	client := auth.SetupClient()
	if glog.IsLevelNormal() {
		util.LogCurrentTrack(client, "playing")
	}
	if err := client.Play(); err != nil {
		glog.Fatal("couldn't start playback: %s", err)
	}
	return nil
}
func mainPause(c *cli.Context) error {
	client := auth.SetupClient()
	if glog.IsLevelNormal() {
		util.LogCurrentTrack(client, "pausing")
	}
	if err := client.Pause(); err != nil {
		glog.Fatal("couldn't pause playback: %s", err)
	}
	return nil
}

func mainNext(c *cli.Context) error {
	client := auth.SetupClient()
	util.LogCurrentTrack(client, "skipping")
	if err := client.Next(); err != nil {
		glog.Fatal("couldn't skip track: ", err)
	}
	return nil
}
func mainPrev(c *cli.Context) error {
	client := auth.SetupClient()
	if err := client.Previous(); err != nil {
		glog.Fatal("couldn't go back: ", err)
	}
	util.LogCurrentTrack(client, "going back to")
	return nil
}

func main() {
	app := cli.NewApp()
	app.Writer = os.Stderr
	app.ErrWriter = os.Stderr
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:     "play",
			Category: "play control",
			Usage:    "play the current song",
			Action:   mainPlay,
		},
		{
			Name:     "pause",
			Category: "play control",
			Usage:    "pause the current song",
			Action:   mainPause,
		},
		{
			Name:     "skip",
			Category: "play control",
			Aliases:  []string{"next"},
			Usage:    "skip the current song",
			Action:   mainNext,
		},
		{
			Name:     "prev",
			Category: "play control",
			Aliases:  []string{"back"},
			Usage:    "prev the current song",
			Action:   mainPrev,
		},
		{
			Name:   "fav",
			Usage:  "add current song to library",
			Action: mainFav,
		},
	}
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
