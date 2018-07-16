package util

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/brianloveswords/spotify/logger"
	"github.com/fatih/color"
	"github.com/lpabon/godbc"
	"github.com/zmb3/spotify"
)

var glog = logger.DefaultLogger

func MustGet(v string, key string) string {
	if envVal := os.Getenv(key); envVal != "" {
		v = envVal
	}
	godbc.Ensure(v != "", fmt.Sprintf("couldn't ensure value with key %s", key))
	return v
}

func MustCreate(filename string) *os.File {
	f, err := os.Create(filename)
	if err != nil {
		panic(err.Error())
	}
	return f
}

func LoadIntMap(name string) (intmap map[string]int) {
	f, err := os.Open(name)
	if err != nil {
		return make(map[string]int)
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&intmap); err != nil {
		panic(err.Error())
	}
	return intmap
}

func SaveIntMap(name string, intmap map[string]int) {
	f := MustCreate(name)
	defer f.Close()

	enc := gob.NewEncoder(f)
	if err := enc.Encode(intmap); err != nil {
		panic(err)
	}
}

func OpenURL(url string, fallback bool) {
	cmd := exec.Command("open", url)
	if err := cmd.Run(); err != nil && fallback {
		// fall back to just printing it
		fmt.Printf("go here\n: %s\n", url)
		return
	}
}

func SongAttributionFromTrack(track *spotify.FullTrack) string {
	song := track.Name
	return fmt.Sprintf("%s - %s", ArtistFromTrack(track), song)
}

func ArtistFromTrack(track *spotify.FullTrack) string {
	var artists []string
	for _, artist := range track.Artists {
		artists = append(artists, artist.Name)
	}
	return strings.Join(artists, ", ")
}
func SongAttributionFromSimpleTrack(track *spotify.SimpleTrack) string {
	return SongAttributionFromTrack(&spotify.FullTrack{
		SimpleTrack: *track,
	})
}

func ArtistFromSimpleTrack(track *spotify.SimpleTrack) string {
	return ArtistFromTrack(&spotify.FullTrack{
		SimpleTrack: *track,
	})
}

func TracksToIDs(tracks []spotify.SimpleTrack) (ids []spotify.ID) {
	for _, track := range tracks {
		ids = append(ids, track.ID)
	}
	return ids
}

func MustGetCurrentlyPlaying(client *spotify.Client, glog logger.Logger) *spotify.FullTrack {
	playing, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		glog.Fatal("could not get currently playing: %s", err)
	}
	return playing.Item
}

func RandomTracks(tracks []spotify.SimpleTrack, n int) (results []spotify.SimpleTrack) {
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

func GetAllTracksByArtist(client *spotify.Client, artistID spotify.ID) (alltracks []spotify.SimpleTrack, err error) {
	defer glog.Enter("util.GetAllTracksByArtist")()

	albums, err := GetAllAlbumsByArtist(client, artistID)
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
func GetAllAlbumsByArtist(client *spotify.Client, artistID spotify.ID) ([]spotify.SimpleAlbum, error) {
	defer glog.Enter("util.GetAllAlbumsByArtist")()
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

func FindArtistID(c *spotify.Client, artist string) *spotify.ID {
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

func LogCurrentTrack(client *spotify.Client, glog logger.Logger, prefix string) {
	playing, _ := client.PlayerCurrentlyPlaying()
	if playing != nil {
		song := SongAttributionFromTrack(playing.Item)
		glog.Log("%s %s", prefix, color.CyanString(song))
	}
}
