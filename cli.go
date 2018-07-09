package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/zmb3/spotify"
)

var (
	paramMixtapeArtist   *string
	paramMixtapeTrackID  *spotify.ID
	paramMixtapeLength   int
	defaultMixtapeLength = 10
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

// func usageAndExit() {
// 	// TODO: fill this out
// 	glog.Fatal("spotify usage goes here")
// }

func cliRouter(args []string) {
	setLogLevel()

	var (
		mainCmd    Command
		mixtapeCmd Command
		showCmd    Command
	)

	mainCmd = Command{
		Name: "spotify",
		Help: "control spotify from the commandline",
		Commands: Subcommands{
			&Command{
				Name: "play",
				Help: "play the current song",
				Func: mainPlay,
			},
			&Command{
				Name: "pause",
				Help: "pause the current song",
				Func: mainPause,
			},
			&Command{
				Name:  "skip",
				Help:  "skip the current song",
				Alias: []string{"next"},
				Func:  mainNext,
			},
			&Command{
				Name:  "prev",
				Help:  "go back to the last song",
				Alias: []string{"back"},
				Func:  mainPrev,
			},
			&Command{
				Name: "fav",
				Help: "add the current song to your library",
				Func: mainFav,
			},
			&showCmd,
			&mixtapeCmd,
		},
	}

	showCmd = Command{
		Name: "show",
		Help: "show the current playing track. has subcommands",
		Func: showShow,
		Commands: Subcommands{
			&Command{
				Name: "artist",
				Help: "show the artist of the current track",
				Func: showArtist,
			},
			&Command{
				Name: "artist-id",
				Help: "show the spotify artist ID of current track",
				Func: showArtistID,
			},
			&Command{
				Name: "artist-uri",
				Help: "show the spotify URI for the artist of current track",
				Func: showArtistURI,
			},
			&Command{
				Name: "track",
				Help: "show the title for the current track",
				Func: showTrack,
			},
			&Command{
				Name: "track-id",
				Help: "show the track ID of the current track",
				Func: showTrackID,
			},
			&Command{
				Name: "track-uri",
				Help: "show the spotify URI for the current track",
				Func: showTrackURI,
			},
		},
	}

	mixtapeCmd = Command{
		Name: "mixtape",
		Help: "create a mixtape",
		Func: mixtapeCreate,
		Examples: []Example{
			Example{
				Args: []string{`artist`},
				Desc: "create a mixtape from artist currently playing",
			},
			Example{
				Args: []string{`artist`, `length=20`},
				Desc: "create a mixtape of 20 tracks from artist currently playing",
			},
			Example{
				Args: []string{`artist="The Sword"`},
				Desc: "create a mixtape of songs by The Sword",
			},
			Example{
				Args: []string{`artist=bill`},
				Desc: "if the artist search is ambiguous, you will be prompted with selections",
			},
			Example{
				Args: []string{`track`},
				Desc: "create mixtape of tracks recommended based on current playing track",
			},
			Example{
				Negative: true,
				Args:     []string{`artist=Chavez`, `track`},
				Comment:  "THIS WON'T WORK!",
				Desc:     "note `artist` and `track` are mutually exclusive, don't include both",
			},
		},
		Params: []Param{
			Param{
				Name:  "artist",
				Alias: []string{"a"},
				Help:  "artist to base mixtape on. pass empty artist to use current playing",
				ParseFn: func(val string) error {
					if paramMixtapeTrackID != nil {
						return fmt.Errorf("must pass track or artist, but not both")
					}
					paramMixtapeArtist = &val
					return nil
				},
			},
			Param{
				Name:  "track",
				Alias: []string{"t"},
				Help:  "track ID to base mixtape on. pass empty track to use current playing",
				ParseFn: func(val string) error {
					if paramMixtapeArtist != nil {
						return fmt.Errorf("must pass track or artist, but not both")
					}
					trackID := spotify.ID(val)
					paramMixtapeTrackID = &trackID
					return nil
				},
			},
			Param{
				Name:     "length",
				Alias:    []string{"n"},
				Help:     fmt.Sprintf("number of tracks to include. defaults to %d", defaultMixtapeLength),
				Implicit: true,
				ParseFn: func(val string) (err error) {
					if val == "" {
						paramMixtapeLength = defaultMixtapeLength
						return nil
					}
					paramMixtapeLength, err = strconv.Atoi(val)
					if paramMixtapeLength < 1 || paramMixtapeLength > 50 {
						return fmt.Errorf("mixtape length must be between 1 and 50")
					}
					return err
				},
			},
		},
	}

	err := mainCmd.Run(args)
	if err != nil {
		glog.Fatal("%s", err)
	}
}

func mustGetCurrentlyPlaying() *spotify.FullTrack {
	client := setupClient()
	playing, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		glog.Fatal("could not get currently playing: %s", err)
	}
	return playing.Item
}

// func playlistCreate(args []string) {
// 	defer glog.Enter("playlistCreate")()

// 	if len(args) == 0 {
// 		glog.Log("err: not enough arguments found for `create`")
// 		usageAndExit()
// 	}
// 	switch playlistCreateParse(args) {
// 	case "songkick-show":
// 		playlistFromSongkickShowPage(args[0])
// 	case "plain":
// 		glog.Log("TODO: create plain playlist")
// 	}
// }

var reURL = regexp.MustCompile("^https?://")

func createPlaylist(client *spotify.Client, name string) *spotify.FullPlaylist {
	user, err := client.CurrentUser()
	if err != nil {
		glog.Fatal("couldn't access current user: %s", err)
	}

	playlist, err := client.CreatePlaylistForUser(user.ID, name, true)
	if err != nil {
		glog.Fatal("couldn't create playlist for user %s: %s", user.ID, err)
	}

	return playlist
}

func getAllTracksByArtist(client *spotify.Client, artistID spotify.ID) (alltracks []spotify.SimpleTrack, err error) {
	defer glog.Enter("getAllTracksByArtist")()

	albums, err := getAllAlbumsByArtist(client, artistID)
	if err != nil {
		return nil, err
	}

	for _, album := range albums {
		page, err := client.GetAlbumTracks(album.ID)
		if err != nil {
			glog.Log("couldn't get tracks for %s (%s): %s", album.Name, album.ID, err)
			continue
		}

		for _, track := range page.Tracks {
			// an album that's attributed to an artist might be a split,
			// so we don't want to add all the songs on the record, just
			// the ones by the artist we're lookin for
			for _, artist := range track.Artists {
				if artist.ID == artistID {
					alltracks = append(alltracks, track)
				}
			}
		}

	}

	return alltracks, nil
}

func getAllAlbumsByArtist(client *spotify.Client, artistID spotify.ID) ([]spotify.SimpleAlbum, error) {
	// TODO: ensure artistID looks like an artistID

	// TODO: some artists may have more than 50 albums but fuck them
	limit := 50
	// TODO: limit to singles and albums or else a lot more artists are
	// going to get more than 50 results and we don't wanna deal with
	// that right now
	albumType := spotify.AlbumTypeSingle | spotify.AlbumTypeAlbum
	page, err := client.GetArtistAlbumsOpt(artistID, &spotify.Options{
		Limit: &limit,
	}, &albumType)
	if err != nil {
		glog.Debug("error getting albums for artist by id %s", artistID)
		return nil, err
	}
	return page.Albums, nil
}

func randomTracks(tracks []spotify.SimpleTrack, n int) (results []spotify.SimpleTrack) {
	defer glog.Enter("randomTracks")()
	max := len(tracks)

	// if there are less tracks than we want to grab, we don't need to
	// do any work, just return it
	if max < n {
		return tracks
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// we want to keep track of what tracks we've seen so we don't end
	// up with a playlist that has duplicates
	seen := make(map[spotify.ID]bool)
	for i := 0; len(results) < n; i++ {
		track := tracks[r.Intn(max)]
		if !seen[track.ID] {
			seen[track.ID] = true
			results = append(results, track)
		}
	}
	return results
}

func playlistCreateParse(args []string) string {
	defer glog.Enter("playlistCreateParse")()

	input := args[0]
	if reURL.MatchString(input) {
		if strings.HasPrefix(input, "https://www.songkick.com/concerts/") {
			return "songkick-show"
		}
		glog.Fatal("don't know what to do with url %s", input)
	}

	if args[0] == "random-by-artist-id" {
		if len(args[1:]) != 1 {
			glog.Fatal("err: playlist create random-by-artist-id expects 1 argument")
		}
		return "random-by-artist-id"
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

func mainFav() {
	defer glog.Enter("mainFav")()
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

func mainNext() {
	defer glog.Enter("mainNext")()
	client := setupClient()
	logCurrentTrack(client, "skipping")
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
	logCurrentTrack(client, "going back to")
}

func logCurrentTrack(client *spotify.Client, prefix string) {
	playing, _ := client.PlayerCurrentlyPlaying()
	if playing != nil {
		song := songAttributionFromTrack(playing.Item)
		glog.Log("%s %s", prefix, color.CyanString(song))
	}
}
