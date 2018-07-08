package main

import (
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/zmb3/spotify"
)

var glog = NewLogger()

func setLogLevel() {
	switch os.Getenv("LOGLEVEL") {
	case "silent":
		glog.Level = LoggerLevelSilent

	case "":
		fallthrough
	case "normal":
		glog.Level = LoggerLevelNormal

	case "verbose":
		glog.Level = LoggerLevelVerbose

	case "debug":
		glog.Level = LoggerLevelDebug

	case "extreme":
		glog.Level = LoggerLevelExtreme
	}
}

func usageAndExit() {
	// TODO: fill this out
	glog.Fatal("spotify usage goes here")
}

func cliRouter(args []string) {
	setLogLevel()

	if len(args) == 0 {
		usageAndExit()
	}

	switch subcmd := args[0]; subcmd {
	case "play":
		mainPlay()

	case "pause":
		mainPause()

	case "skip":
		fallthrough
	case "next":
		mainNext()

	case "prev":
		fallthrough
	case "previous":
		mainPrev()

	case "playing":
		playingRouter(args[1:])

	case "playlist":
		playlistRouter(args[1:])

	default:
		glog.Log("err: %s not a valid subcommand", subcmd)
		usageAndExit()
	}
}

func playingRouter(args []string) {
	if len(args) == 0 {
		glog.Log("err: missing command for playing")
		usageAndExit()
	}

	switch subcmd := args[0]; subcmd {
	case "fav":
		playingFav()
	case "show":
		playingShow(args[1:])
	default:
		glog.Log("err: %s not a valid subcommand", subcmd)
		usageAndExit()
	}
}

func playlistRouter(args []string) {
	if len(args) == 0 {
		glog.Log("err: missing command for playlist")
		usageAndExit()
	}

	switch subcmd := args[0]; subcmd {
	case "create":
		playlistCreate(args[1:])
	default:
		glog.Log("err: %s not a valid subcommand", subcmd)
		usageAndExit()
	}
}

func playlistCreate(args []string) {
	defer glog.Enter("playlistCreate")()

	if len(args) == 0 {
		glog.Log("err: not enough arguments found for `create`")
		usageAndExit()
	}
	switch playlistCreateParse(args) {
	case "songkick-show":
		playlistFromSongkickShowPage(args[0])
		// create playlist from songkick show page
	case "plain":
		glog.Log("TODO: create plain playlist")
		// create playlist by the name given
	}
}

var reURL = regexp.MustCompile("^https?://")

func playlistCreateParse(args []string) string {
	defer glog.Enter("playlistCreateParse")()

	input := args[0]
	if reURL.MatchString(input) {
		if strings.HasPrefix(input, "https://www.songkick.com/concerts/") {
			return "songkick-show"
		}
		glog.Fatal("don't know what to do with url %s", input)
	}

	return "plain"
}

func playlistFromSongkickShowPage(url string) {
	defer glog.Enter("playlistFromSongkickShowPage")()
	glog.Verbose("creating playlist from page %s", color.YellowString(url))

	client := setupClient()
	artists := artistsFromSongkickShowPage(url)
	name := strings.Join(artists, "/")

	glog.Log("creating playlist %s", color.CyanString(name))

	playlist := createPlaylist(client, name)
	addArtistLatestAlbumsPlaylist(client, playlist, artists)

	// TODO: option to open the resulting spotify playlist??
	glog.CmdOutput("%s", playlist.URI)
}

func addArtistLatestAlbumsPlaylist(
	client *spotify.Client,
	playlist *spotify.FullPlaylist,
	artists []string,
) *spotify.FullPlaylist {
	defer glog.Enter("addArtistLatestAlbumsPlaylist")()
	for _, artist := range artists {
		id := findArtistID(client, artist)
		if id == nil {
			glog.Log("couldn't find an artist result for %s", color.RedString(artist))
			continue
		}
		albums := getLatestAlbums(client, *id)
		addAlbumsToPlaylist(client, playlist, albums)
	}
	return playlist
}

func playingShow(args []string) {
	defer glog.Enter("playingFav")()
	client := setupClient()
	playing, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		glog.Fatal("could not get currently playing: %s", err)
	}
	track := playing.Item

	if len(args) == 0 {
		name := songAttributionFromTrack(track)
		glog.Log("currently playing %s", color.CyanString(name))
		return
	}
	switch args[0] {
	case "help":
		fallthrough
	case "--help":
		glog.Fatal("TODO: implement help")
	case "artist":
		glog.CmdOutput("%s", track.Artists[0].Name)
	case "artist-id":
		glog.CmdOutput("%s", track.Artists[0].ID)
	case "artist-uri":
		glog.CmdOutput("%s", track.Artists[0].URI)
	case "track":
		glog.CmdOutput("%s", track.Name)
	case "track-id":
		glog.CmdOutput("%s", track.ID)
	case "track-uri":
		glog.CmdOutput("%s", track.URI)
	default:
		glog.Debug("fell through switch")
		glog.Fatal("I don't understand argument %s", color.RedString(args[0]))
	}
}

func playingFav() {
	defer glog.Enter("playingFav")()
	client := setupClient()
	playing, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		glog.Fatal("could not get currently playing: %s", err)
	}
	track := playing.Item
	if err := client.AddTracksToLibrary(track.ID); err != nil {
		glog.Fatal("could add track to library: %s", err)
	}
	name := songAttributionFromTrack(track)
	glog.Log("adding to library: %s", color.CyanString(name))
}

func mainPlay() {
	defer glog.Enter("mainPlay")()
	client := setupClient()
	if glog.IsLevelNormal() {
		logCurrentTrack(client, "playing")
	}
	if err := client.Play(); err != nil {
		glog.Fatal("couldn't start playback: %s", err)
	}
}
func mainPause() {
	defer glog.Enter("mainPause")()
	client := setupClient()
	if glog.IsLevelNormal() {
		logCurrentTrack(client, "pausing")
	}
	if err := client.Pause(); err != nil {
		glog.Fatal("couldn't pause playback: %s", err)
	}
}

func logCurrentTrack(client *spotify.Client, prefix string) {
	playing, _ := client.PlayerCurrentlyPlaying()
	if playing != nil {
		song := songAttributionFromTrack(playing.Item)
		glog.Log("%s %s", prefix, color.CyanString(song))
	}
}

func mainNext() {
	defer glog.Enter("mainNext")()
	client := setupClient()
	if err := client.Next(); err != nil {
		glog.Fatal("couldn't skip track: ", err)
	}
}
func mainPrev() {
	defer glog.Enter("mainPrev")()
	client := setupClient()
	if err := client.Previous(); err != nil {
		glog.Fatal("couldn't go back: ", err)
	}
}
