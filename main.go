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
	cliRouter(os.Args[1:])
}

func findArtistID(c *spotify.Client, artist string) *spotify.ID {
	page, err := c.Search(artist, spotify.SearchTypeArtist)
	if err != nil {
		panic(err)
	}

	artists := page.Artists.Artists

	if len(artists) == 1 {
		return &artists[0].ID
	}

	for _, found := range artists {
		if found.Name == artist {
			return &found.ID
		}
	}
	return nil
}

func getAlbumTracks(c *spotify.Client, album *spotify.SimpleAlbum) []spotify.SimpleTrack {
	tracks, err := c.GetAlbumTracks(album.ID)
	if err != nil {
		panic(err)
	}
	return tracks.Tracks
}

func tracksToIDs(tracks []spotify.SimpleTrack) (ids []spotify.ID) {
	for _, track := range tracks {
		ids = append(ids, track.ID)
	}
	return ids
}

func addAlbumsToPlaylist(
	client *spotify.Client,
	playlist *spotify.FullPlaylist,
	albums []spotify.SimpleAlbum,
) {
	user, _ := client.CurrentUser()
	for _, album := range albums {
		tracks := getAlbumTracks(client, &album)
		_, err := client.AddTracksToPlaylist(user.ID, playlist.ID, tracksToIDs(tracks)...)
		if err != nil {
			panic(err)
		}
	}
}

// we want to get all of the latest singles released after the latest
// album, and also the latest album
func getLatestAlbums(client *spotify.Client, artistID spotify.ID) []spotify.SimpleAlbum {
	limit := 50
	albumType := spotify.AlbumTypeSingle | spotify.AlbumTypeAlbum
	page, err := client.GetArtistAlbumsOpt(artistID, &spotify.Options{
		Limit: &limit,
	}, &albumType)

	albums := page.Albums

	if err != nil {
		panic(err)
	}

	// TODO: assumes reverse chronological order, make sure that ends up
	// holding to be generally true
	var results []spotify.SimpleAlbum
	for _, album := range albums {
		results = append(results, album)
		if album.AlbumType == "album" {
			break
		}
	}

	return results
}

func getCurrentArtistID(client *spotify.Client) spotify.ID {
	playing, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		panic(err)
	}
	artist := artistFromTrack(playing.Item)
	id := playing.Item.Artists[0].ID
	debugprint("%s %s\n", artist, id)
	return id
}

func createPlaylist(client *spotify.Client, name string) *spotify.FullPlaylist {
	user, _ := client.CurrentUser()
	playlist, err := client.CreatePlaylistForUser(user.ID, name, true)
	if err != nil {
		panic(err)
	}
	return playlist
}

func songAttributionFromTrack(track *spotify.FullTrack) string {
	song := track.Name
	return fmt.Sprintf("%s - %s", artistFromTrack(track), song)
}

func artistFromTrack(track *spotify.FullTrack) string {
	var artists []string
	for _, artist := range track.Artists {
		artists = append(artists, artist.Name)
	}
	return strings.Join(artists, ", ")
}
