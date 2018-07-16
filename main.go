package main

import (
	"os"

	"github.com/brianloveswords/spotify/auth"
	"github.com/brianloveswords/spotify/logger"
	"github.com/brianloveswords/spotify/mix"
	"github.com/brianloveswords/spotify/util"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"github.com/zmb3/spotify"
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
	if err := client.Play(); err != nil {
		glog.Fatal("couldn't start playback: %s", err)
	}
	if glog.IsLevelNormal() {
		util.LogCurrentTrack(client, glog, "playing")
	}
	return nil
}
func mainPause(c *cli.Context) error {
	client := auth.SetupClient()
	if err := client.Pause(); err != nil {
		glog.Fatal("couldn't pause playback: %s", err)
	}
	if glog.IsLevelNormal() {
		util.LogCurrentTrack(client, glog, "pausing")
	}
	return nil
}

func mainNext(c *cli.Context) error {
	client := auth.SetupClient()
	if err := client.Next(); err != nil {
		glog.Fatal("couldn't skip track: ", err)
	}
	util.LogCurrentTrack(client, glog, "skipping")
	return nil
}
func mainPrev(c *cli.Context) error {
	client := auth.SetupClient()
	if err := client.Previous(); err != nil {
		glog.Fatal("couldn't go back: ", err)
	}
	util.LogCurrentTrack(client, glog, "going back to")
	return nil
}

func showPlaying(c *cli.Context) error {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient(), glog)
	name := util.SongAttributionFromTrack(track)
	glog.CmdOutput("current track: %s", color.CyanString(name))
	return nil
}
func showArtist(c *cli.Context) error {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient(), glog)
	glog.CmdOutput(track.Artists[0].Name)
	return nil
}
func showArtistID(c *cli.Context) error {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient(), glog)
	glog.CmdOutput("%s", track.Artists[0].ID)
	return nil
}
func showArtistURI(c *cli.Context) error {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient(), glog)
	uri := track.Artists[0].URI
	glog.CmdOutput("%s", uri)

	if c.Bool("open") {
		util.OpenURL(string(uri), false)
	}

	return nil
}
func showTrack(c *cli.Context) error {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient(), glog)
	glog.CmdOutput("%s", track.Name)
	return nil
}
func showTrackID(c *cli.Context) error {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient(), glog)
	glog.CmdOutput("%s", track.ID)
	return nil
}
func showTrackURI(c *cli.Context) error {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient(), glog)
	glog.CmdOutput("%s", track.URI)

	if c.Bool("open") {
		util.OpenURL(string(track.URI), false)
	}
	return nil
}
func showAlbum(c *cli.Context) error {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient(), glog)
	glog.CmdOutput("%s", track.Album.Name)
	return nil
}
func showAlbumID(c *cli.Context) error {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient(), glog)
	glog.CmdOutput("%s", track.Album.ID)
	return nil
}
func showAlbumURI(c *cli.Context) error {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient(), glog)
	glog.CmdOutput("%s", track.Album.URI)

	if c.Bool("open") {
		util.OpenURL(string(track.Album.URI), false)
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Writer = &glog
	app.ErrWriter = &glog
	app.Version = "1.0.0"

	flagMixLength := cli.IntFlag{
		Name:  "length, l",
		Usage: "how many tracks to include in the mix",
		Value: 10,
	}

	flagOpen := cli.BoolFlag{
		Name:  "open",
		Usage: "opens in spotify, if possible",
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "silent",
			Usage: "don't output any logging and don't wait for stdin",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "output additional information while running commands",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "output debugging information while running commands",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool("verbose") {
			glog.Level = logger.LevelVerbose
		}
		if c.Bool("debug") {
			glog.Level = logger.LevelDebug
		}
		if c.Bool("silent") {
			glog.Level = logger.LevelSilent
		}
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:   "fav",
			Usage:  "add current song to library",
			Action: mainFav,
		},
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
			Name:     "info",
			Category: "track information",
			Usage:    "display current song",
			Action:   showPlaying,
		},
		{
			Name:     "track",
			Category: "track information",
			Usage:    "show the name of the current song",
			Action:   showTrack,
		},
		{
			Name:     "track-id",
			Category: "track information",
			Usage:    "show the ID of the current song",
			Action:   showTrackID,
		},
		{
			Name:     "track-uri",
			Category: "track information",
			Usage:    "show the URI of the current song",
			Action:   showTrackURI,
			Flags:    []cli.Flag{flagOpen},
		},
		{
			Name:     "artist",
			Category: "track information",
			Usage:    "show the name of the current artist",
			Action:   showArtist,
		},
		{
			Name:     "artist-id",
			Category: "track information",
			Usage:    "show the ID of the current artist",
			Action:   showArtistID,
		},
		{
			Name:     "artist-uri",
			Category: "track information",
			Usage:    "show the spotify URI of the current artist",
			Action:   showArtistURI,
			Flags:    []cli.Flag{flagOpen},
		},
		{
			Name:     "album",
			Category: "track information",
			Usage:    "show the name of the current album",
			Action:   showAlbum,
		},
		{
			Name:     "album-id",
			Category: "track information",
			Usage:    "show the ID of the current album",
			Action:   showAlbumID,
		},
		{
			Name:     "album-uri",
			Category: "track information",
			Usage:    "show the spotify URI of the current album",
			Action:   showAlbumURI,
			Flags:    []cli.Flag{flagOpen},
		},
		{
			Name:  "mix",
			Usage: "commands for creating mixes",
			Subcommands: []cli.Command{
				{
					Name:      "track",
					Usage:     "create mix from track. if no track given, uses current track",
					ArgsUsage: "[trackID]",
					Action:    mixTrack,
					Flags: []cli.Flag{
						flagOpen,
						flagMixLength,
						cli.StringFlag{
							Name:  "name",
							Usage: "what to call the playlist",
							Value: "{mix} :ARTIST: - :TRACK:",
						},
					},
				},
				{
					Name:      "artist",
					Usage:     "create mix from artist.",
					UsageText: "if no artist given, uses current artist. if artist is ambiguous, gives options\n   to select from on stdin, unless --silent",
					ArgsUsage: "[artist]",
					Action:    mixArtist,
					Flags: []cli.Flag{
						flagOpen,
						flagMixLength,
						cli.BoolFlag{
							Name:  "id",
							Usage: "interpret argument as artist ID",
						},
						cli.StringFlag{
							Name:  "name",
							Usage: "what to call the playlist",
							Value: "{mix} :ARTIST:",
						},
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func mixTrack(c *cli.Context) error {
	defer glog.Enter("mixTrack")()
	var (
		playlist *spotify.FullPlaylist
		err      error
		length   = c.Int("length")
		track    = c.Args().Get(0)
		name     = c.String("name")
	)

	glog.Debug("name %q", name)
	glog.Debug("length %q", length)

	if track == "" {
		playlist, err = mix.ByCurrentTrack(glog, name, length)
	} else {
		playlist, err = mix.ByTrackID(glog, spotify.ID(track), name, length)
	}
	if err != nil {
		glog.Fatal(err.Error())
	}

	glog.Log("created %s", color.MagentaString(playlist.Name))
	glog.CmdOutput("%s", playlist.URI)

	if c.Bool("open") {
		util.OpenURL(string(playlist.URI), false)
	}
	return nil
}

func mixArtist(c *cli.Context) error {
	defer glog.Enter("mixTrack")()
	var (
		playlist *spotify.FullPlaylist
		err      error
		artist   = c.Args().Get(0)
		length   = c.Int("length")
		name     = c.String("name")
		isID     = c.Bool("id")
	)

	glog.Debug("artist %q", name)
	glog.Debug("name %q", name)
	glog.Debug("length %q", length)

	if artist == "" {
		if isID {
			glog.Fatal("must pass an artist ID when using --id flag")
		}
		playlist, err = mix.ByCurrentArtist(glog, name, length)
	} else {
		if isID {
			playlist, err = mix.ByArtistID(glog, spotify.ID(artist), name, length)
		} else {
			playlist, err = mix.ByArtist(glog, artist, name, length)
		}
	}
	if err != nil {
		glog.Fatal(err.Error())
	}

	glog.Log("created %s", color.MagentaString(playlist.Name))
	glog.CmdOutput("%s", playlist.URI)

	if c.Bool("open") {
		util.OpenURL(string(playlist.URI), false)
	}
	return nil
}
