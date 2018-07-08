package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/zmb3/spotify"
)

func usageAndExit() {
	// TODO: fill this out
	fmt.Fprintf(os.Stderr, "spotify usage goes here\n")
	os.Exit(1)
}

func cliRouter(args []string) {
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
		fmt.Fprintf(os.Stderr, "err: %s not a valid subcommand\n", subcmd)
		usageAndExit()
	}
}

func playingRouter(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "err: missing command for playing\n")
		usageAndExit()
	}

	switch subcmd := args[0]; subcmd {
	case "fav":
		playingFav()
	default:
		fmt.Fprintf(os.Stderr, "err: %s not a valid subcommand\n", subcmd)
		usageAndExit()
	}
}

func playlistRouter(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "err: missing command for playlist\n")
		usageAndExit()
	}

	switch subcmd := args[0]; subcmd {
	case "create":
		playlistCreate(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "err: %s not a valid subcommand\n", subcmd)
		usageAndExit()
	}
}

func playlistCreate(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "err: not enough arguments found for `create`\n")
		usageAndExit()
	}
	switch playlistCreateParse(args) {
	case "songkick-show":
		fmt.Println("creating playlist from songkick show page")
		playlistFromSongkickShowPage(args[0])
		// create playlist from songkick show page
	case "plain":
		fmt.Println("creating plain playlist")
		// create playlist by the name given
	}
}

var reURL = regexp.MustCompile("^https?://")

func playlistCreateParse(args []string) string {
	input := args[0]
	if reURL.MatchString(input) {
		if strings.HasPrefix(input, "https://www.songkick.com/concerts/") {
			return "songkick-show"
		}
		log.Fatalf("don't know what to do with url %s", input)
	}

	return "plain"
}

func playlistFromSongkickShowPage(url string) {
	// TODO: option to open the resulting spotify playlist??

	client := setupClient()
	artists := artistsFromSongkickShowPage(url)
	name := strings.Join(artists, "/")

	fmt.Print("creating playlist ")
	color.Cyan(name)

	playlist := createPlaylist(client, name)
	addArtistLatestAlbumsPlaylist(client, playlist, artists)

	fmt.Println(playlist.URI)
}

func addArtistLatestAlbumsPlaylist(
	client *spotify.Client,
	playlist *spotify.FullPlaylist,
	artists []string,
) *spotify.FullPlaylist {
	for _, artist := range artists {

		id := findArtistID(client, artist)
		if id == nil {
			fmt.Fprintf(os.Stderr, "couldn't find an artist result for %s\n", artist)
			continue
		}

		albums := getLatestAlbums(client, *id)
		addAlbumsToPlaylist(client, playlist, albums)
	}
	return playlist
}

func playingFav() {
	client := setupClient()
	track := addCurrentlyPlayingToLibrary(client)

	fmt.Print("added to library: ")
	color.Cyan(songAttributionFromTrack(track))
}

func mainPlay() {
	client := setupClient()
	if err := client.Play(); err != nil {
		log.Fatal("couldn't start playback: ", err)
	}
}
func mainPause() {
	client := setupClient()
	if err := client.Pause(); err != nil {
		log.Fatal("couldn't pause playback: ", err)
	}
}
func mainNext() {
	client := setupClient()
	if err := client.Next(); err != nil {
		log.Fatal("couldn't skip track: ", err)
	}
}
func mainPrev() {
	client := setupClient()
	if err := client.Previous(); err != nil {
		log.Fatal("couldn't go back: ", err)
	}
}

// client := setupClient()
// user, _ := client.CurrentUser()
// fmt.Fprintf(os.Stderr, "user: %s\n", user.ID)

// tracks := getAllTracks(client)
// artists := processTracklist(tracks)
// printHistogram(artists)
// lookupSongkickIDs(artists)
// addCurrentlyPlayingToLibrary(client)
// fmt.Println(getCurrentArtistID(client))

// page := "https://www.songkick.com/concerts/33692814-royal-they-at-alphaville"
// page := os.Args[1]
// if page == "" {
// 	fmt.Printf("must provide a songkick page")
// 	os.Exit(1)
// }

// artists := artistsFromSongkickPage(page)
// if len(artists) > 0 {
// 	createShowPlaylist(client, artists)
// }
// fmt.Println(artists)
