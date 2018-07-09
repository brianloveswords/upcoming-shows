package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/lpabon/godbc"
	"github.com/zmb3/spotify"
)

var trackDataFilename = "saved-tracks.data"

func loadTrackData() (tracks []spotify.SavedTrack) {
	f, err := os.Open(trackDataFilename)
	if err != nil {
		glog.Debug("couldn't open saved track data from file '%s'\n", trackDataFilename)
		return nil
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&tracks); err != nil {
		panic(err.Error())
	}
	godbc.Ensure(len(tracks) > 0, "should have found at least one track")
	return tracks
}

func saveTrackData(tracks []spotify.SavedTrack) {
	f := mustCreate(trackDataFilename)
	defer f.Close()

	enc := gob.NewEncoder(f)
	if err := enc.Encode(tracks); err != nil {
		panic(err)
	}
}

func getAllTracks(client *spotify.Client) (tracks []spotify.SavedTrack) {
	if savedTracks := loadTrackData(); len(savedTracks) > 0 {
		glog.Debug("loading tracks from disk\n")
		return savedTracks
	}

	var offset int
	limit := 50
	for i := 0; ; i++ {
		offset = limit * i
		page, err := client.CurrentUsersTracksOpt(&spotify.Options{
			Limit:  &limit,
			Offset: &offset,
		})
		if err != nil {
			log.Fatalf("error getting tracks: %v", err)
		}

		glog.Debug("got %s", page.Endpoint)

		tracks = append(tracks, page.Tracks...)

		if len(tracks) >= page.Total {
			break
		}
	}
	saveTrackData(tracks)
	return tracks
}

func processTracklist(tracks []spotify.SavedTrack) (artists []Artist) {
	hist := artistHistogram(tracks)
	for k, v := range hist {
		artists = append(artists, Artist{
			Appearances: v,
			Name:        k,
			SongkickID:  SongkickUnknown,
		})
	}

	sort.Slice(artists, func(i, j int) bool {
		return artists[j].Appearances < artists[i].Appearances
	})
	return artists
}

func printHistogram(artists []Artist) {
	for _, artist := range artists {
		fmt.Printf("%v %s\n", artist.Appearances, artist.Name)
	}
}

func artistHistogram(tracks []spotify.SavedTrack) map[string]int {
	hist := make(map[string]int)
	for _, track := range tracks {
		for _, artist := range track.Artists {
			hist[artist.Name]++
		}
	}
	return hist
}
