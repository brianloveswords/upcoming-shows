package mix

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/brianloveswords/spotify/auth"
	"github.com/brianloveswords/spotify/logger"
	"github.com/brianloveswords/spotify/util"
	"github.com/fatih/color"
	"github.com/zmb3/spotify"
)

var glog = logger.DefaultLogger

func ByCurrentTrack(glog logger.Logger, name string, length int) (*spotify.FullPlaylist, error) {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient(), glog)
	return ByTrackID(glog, track.ID, name, length)
}

func processName(name string, artist *spotify.SimpleArtist, track *spotify.SimpleTrack) string {
	if artist != nil {
		name = strings.Replace(name, ":ARTIST:", artist.Name, -1)
	}
	if track != nil {
		name = strings.Replace(name, ":TRACK:", track.Name, -1)
	}
	return name
}

func ByTrackID(glog logger.Logger, trackID spotify.ID, name string, length int) (*spotify.FullPlaylist, error) {
	client := auth.SetupClient()

	seedTrack, err := client.GetTrack(trackID)
	if err != nil {
		return nil, fmt.Errorf("couldn't find track for trackID %s: %s", trackID, err)
	}
	trackname := util.SongAttributionFromTrack(seedTrack)
	playlistName := processName(name, &seedTrack.Artists[0], &seedTrack.SimpleTrack)
	glog.Log("making mixtape with seed %s...", color.YellowString(trackname))

	seeds := spotify.Seeds{
		Tracks: []spotify.ID{trackID},
	}

	recommendations, err := client.GetRecommendations(seeds, spotify.NewTrackAttributes(), &spotify.Options{
		Limit: &length,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't get recommendations: %s", err)
	}
	tracks := recommendations.Tracks

	user, err := client.CurrentUser()
	if err != nil {
		return nil, fmt.Errorf("couldn't access current user: %s", err)
	}

	playlist, err := client.CreatePlaylistForUser(user.ID, playlistName, true)
	if err != nil {
		return nil, fmt.Errorf("couldn't create playlist for user %s: %s", user.ID, err)
	}

	for _, track := range tracks {
		glog.Verbose("adding %s", color.CyanString(util.SongAttributionFromSimpleTrack(&track)))
	}

	_, err = client.AddTracksToPlaylist(user.ID, playlist.ID, util.TracksToIDs(tracks)...)
	if err != nil {
		// TODO: don't use color formatting here, use structured errors
		return nil, fmt.Errorf("couldn't add tracks to playlist %s for user %s: %s",
			color.BlueString(playlist.Name),
			color.GreenString(user.ID),
			err,
		)
	}
	return playlist, nil
}

func ByArtist(glog logger.Logger, artistName string, name string, length int) (*spotify.FullPlaylist, error) {
	var artistID spotify.ID
	client := auth.SetupClient()
	normalizedArtist := strings.ToLower(artistName)

	page, err := client.Search(artistName, spotify.SearchTypeArtist)
	if err != nil {
		panic(err)
	}

	artists := page.Artists.Artists

	if len(artists) == 0 {
		glog.Fatal("could not find any matches for %s", color.BlueString(artistName))
	}

	if len(artists) == 1 {
		artistID = artists[0].ID
		return ByArtistID(glog, artistID, name, length)
	}

	for _, found := range artists {
		if strings.ToLower(found.Name) == normalizedArtist {
			artistID = found.ID
			return ByArtistID(glog, artistID, name, length)
		}
	}

	// okay we didn't find anything, let's get a user option?
	// if it's silent, we can't prompt, so quit immediately
	if glog.IsLevelSilent() {
		os.Exit(1)
	}

	pick := promptForArtistSelection(artists)

	if pick == nil {
		os.Exit(1)
	}

	return ByArtistID(glog, pick.ID, name, length)
}

func promptForArtistSelection(artists []spotify.FullArtist) *spotify.FullArtist {
	glog.Log("could not find exact match")
	for i, artist := range artists {
		glog.Log("%d) %s", i+1, color.BlueString(artist.Name))
	}
	// read from stdin until we get a valid input
	for {
		glog.Prompt("please select 1-%d", len(artists))
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n ")

		if text == "" {
			glog.Log("giving up")
			return nil
		}

		pick, err := strconv.Atoi(text)
		if err != nil || pick > len(artists) || pick < 1 {
			continue
		}

		return &artists[pick-1]
	}
}
func byArtist(glog logger.Logger, artist spotify.SimpleArtist, name string, length int) (*spotify.FullPlaylist, error) {
	client := auth.SetupClient()
	alltracks, err := util.GetAllTracksByArtist(client, artist.ID)
	if err != nil {
		return nil, fmt.Errorf("could not get tracks from artist with ID %s: %s", artist.ID, err)
	}

	tracks := util.RandomTracks(alltracks, length)
	if len(tracks) == 0 {
		return nil, fmt.Errorf("didn't find any tracks for artist with ID %s", artist.ID)
	}

	user, err := client.CurrentUser()
	if err != nil {
		return nil, fmt.Errorf("couldn't access current user: %s", err)
	}

	playlistName := processName(name, &artist, nil)
	playlist, err := client.CreatePlaylistForUser(user.ID, playlistName, true)
	if err != nil {
		return nil, fmt.Errorf("couldn't create playlist for user %s: %s", user.ID, err)
	}

	for _, track := range tracks {
		glog.Verbose("adding %s", color.CyanString(util.SongAttributionFromSimpleTrack(&track)))
	}

	_, err = client.AddTracksToPlaylist(user.ID, playlist.ID, util.TracksToIDs(tracks)...)
	if err != nil {
		return nil, fmt.Errorf("couldn't add tracks to playlist %s for user %s: %s",
			color.BlueString(playlist.Name),
			color.GreenString(user.ID),
			err,
		)
	}
	return playlist, nil
}

func ByCurrentArtist(glog logger.Logger, name string, length int) (*spotify.FullPlaylist, error) {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient(), glog)
	artist := track.Artists[0]
	return byArtist(glog, artist, name, length)
}

func ByArtistID(glog logger.Logger, artistID spotify.ID, name string, length int) (*spotify.FullPlaylist, error) {
	defer glog.Enter("mixtapeByArtistID")()
	client := auth.SetupClient()

	artist, err := client.GetArtist(artistID)
	if err != nil {
		glog.Fatal("couldn't look up artist with ID %s: %s", artistID, err)
	}

	return byArtist(glog, artist.SimpleArtist, name, length)
}
