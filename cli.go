package main

import (
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

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

	case "mixtape":
		mixtapeRouter(args[1:])

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

func mixtapeRouter(args []string) {
	if len(args) == 0 {
		glog.Log("err: missing command for mixtape")
		usageAndExit()
	}

	switch subcmd := args[0]; subcmd {
	case "current-artist":
		mixtapeByCurrentArtist()
	case "by-artist-id":
		mixtapeByArtistID(args[1])
	default:
		glog.Log("mixtape: %s not a valid subcommand", subcmd)
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
	case "plain":
		glog.Log("TODO: create plain playlist")
	}
}

var reURL = regexp.MustCompile("^https?://")

func mixtapeByCurrentArtist() {
	defer glog.Enter("mixtapeByCurrentArtist")()
	client := setupClient()

	playing, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		glog.Fatal("could not get currently playing: %s", err)
	}
	track := playing.Item
	mixtapeByArtistID(string(track.Artists[0].ID))
}

func mixtapeByArtistID(artistID string) {
	defer glog.Enter("mixtapeByArtistID")()
	client := setupClient()

	artist, err := client.GetArtist(spotify.ID(artistID))
	if err != nil {
		glog.Fatal("couldn't look up artist with ID %s: %s", artistID, err)
	}

	alltracks, err := getAllTracksByArtist(client, artistID)
	if err != nil {
		glog.Fatal("could not get tracks from artist with ID %s: %s", artistID, err)
	}

	tracks := randomTracks(alltracks, 10)
	if len(tracks) == 0 {
		glog.Fatal("didn't find any tracks for artist with ID %s", artistID)
	}

	user, err := client.CurrentUser()
	if err != nil {
		glog.Fatal("couldn't access current user: %s", err)
	}

	playlist, err := client.CreatePlaylistForUser(user.ID, "{mixtape} "+artist.Name, true)
	if err != nil {
		glog.Fatal("couldn't create playlist for user %s: %s", user.ID, err)
	}

	for _, track := range tracks {
		glog.Log("adding %s", color.CyanString(songAttributionFromSimpleTrack(&track)))
	}

	_, err = client.AddTracksToPlaylist(user.ID, playlist.ID, tracksToIDs(tracks)...)
	if err != nil {
		glog.Fatal("couldn't add tracks to playlist %s for user %s: %s",
			color.BlueString(playlist.Name),
			color.GreenString(user.ID),
			err,
		)
	}
	glog.Log("Created playlist %s", color.BlueString(playlist.Name))
	glog.CmdOutput("%s", playlist.URI)
}

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

func getAllTracksByArtist(client *spotify.Client, artistID string) (alltracks []spotify.SimpleTrack, err error) {
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
				if artist.ID == spotify.ID(artistID) {
					alltracks = append(alltracks, track)
				}
			}
		}

	}

	return alltracks, nil
}

func getAllAlbumsByArtist(client *spotify.Client, artistID string) ([]spotify.SimpleAlbum, error) {
	// TODO: ensure artistID looks like an artistID

	// TODO: some artists may have more than 50 albums but fuck them
	limit := 50
	// TODO: limit to singles and albums or else a lot more artists are
	// going to get more than 50 results and we don't wanna deal with
	// that right now
	albumType := spotify.AlbumTypeSingle | spotify.AlbumTypeAlbum
	page, err := client.GetArtistAlbumsOpt(spotify.ID(artistID), &spotify.Options{
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
