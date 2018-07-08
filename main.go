package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/zmb3/spotify"
)

type Artist struct {
	Name        string
	Appearances int
	SongkickID  int
}

func main() {
	client := setupClient()

	if tok, err := client.Token(); err == nil {
		debugprint("token expires: %v", tok.Expiry)
		saveToken(tok, tokenPath)
	}

	user, _ := client.CurrentUser()
	fmt.Fprintf(os.Stderr, "user: %s\n", user.ID)

	// tracks := getAllTracks(client)
	// artists := processTracklist(tracks)
	// printHistogram(artists)
	// lookupSongkickIDs(artists)

	addCurrentlyPlayingToLibrary(client)

}

func addCurrentlyPlayingToLibrary(client *spotify.Client) {
	playing, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		panic(err)
	}

	if err := client.AddTracksToLibrary(playing.Item.ID); err != nil {
		panic(err)
	}

	fmt.Printf("%s (%s) added to library\n",
		songAttributionFromTrack(playing.Item),
		playing.Item.ID)
}

func songAttributionFromTrack(track *spotify.FullTrack) string {
	var artists []string
	song := track.Name
	for _, artist := range track.Artists {
		artists = append(artists, artist.Name)
	}
	artistline := strings.Join(artists, ", ")
	return fmt.Sprintf("%s - %s", artistline, song)
}
