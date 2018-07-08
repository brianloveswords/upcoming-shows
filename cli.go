package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
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
	case "next":
		mainNext()

	case "prev":
	case "previous":
		mainPrev()

	case "playing":
		playingRouter(args[1:])

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
