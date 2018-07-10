package mixtape

import (
	"bufio"
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

func ByCurrentTrack(length int) {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient())
	ByTrackID(track.ID, length)
}

func ByTrackID(trackID spotify.ID, length int) {
	client := auth.SetupClient()

	seedTrack, err := client.GetTrack(trackID)
	if err != nil {
		glog.Fatal("couldn't find track for trackID %s: %s", trackID, err)
	}
	trackname := util.SongAttributionFromTrack(seedTrack)
	glog.Log("making mixtape with seed %s...", color.YellowString(trackname))

	seeds := spotify.Seeds{
		Tracks: []spotify.ID{trackID},
	}
	recommendations, err := client.GetRecommendations(seeds, spotify.NewTrackAttributes(), &spotify.Options{
		Limit: &length,
	})
	if err != nil {
		glog.Fatal("couldn't get recommendations: %s", err)
	}
	tracks := recommendations.Tracks

	user, err := client.CurrentUser()
	if err != nil {
		glog.Fatal("couldn't access current user: %s", err)
	}

	playlist, err := client.CreatePlaylistForUser(user.ID, "{mixtape} "+trackname, true)
	if err != nil {
		glog.Fatal("couldn't create playlist for user %s: %s", user.ID, err)
	}

	for _, track := range tracks {
		glog.Log("adding %s", color.CyanString(util.SongAttributionFromSimpleTrack(&track)))
	}

	_, err = client.AddTracksToPlaylist(user.ID, playlist.ID, util.TracksToIDs(tracks)...)
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

func ByArtist(artist string, length int) {
	var artistID spotify.ID
	client := auth.SetupClient()
	normalizedArtist := strings.ToLower(artist)

	page, err := client.Search(artist, spotify.SearchTypeArtist)
	if err != nil {
		panic(err)
	}

	artists := page.Artists.Artists

	if len(artists) == 0 {
		glog.Fatal("could not find any matches for %s", color.BlueString(artist))
	}

	if len(artists) == 1 {
		artistID = artists[0].ID
		ByArtistID(artistID, length)
		return
	}

	for _, found := range artists {
		if strings.ToLower(found.Name) == normalizedArtist {
			artistID = found.ID
			ByArtistID(artistID, length)
			return
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

	ByArtistID(pick.ID, length)
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

func ByArtistID(artistID spotify.ID, length int) {
	defer glog.Enter("mixtapeByArtistID")()
	client := auth.SetupClient()

	artist, err := client.GetArtist(artistID)
	if err != nil {
		glog.Fatal("couldn't look up artist with ID %s: %s", artistID, err)
	}

	alltracks, err := util.GetAllTracksByArtist(client, artistID)
	if err != nil {
		glog.Fatal("could not get tracks from artist with ID %s: %s", artistID, err)
	}

	tracks := util.RandomTracks(alltracks, length)
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
		glog.Log("adding %s", color.CyanString(util.SongAttributionFromSimpleTrack(&track)))
	}

	_, err = client.AddTracksToPlaylist(user.ID, playlist.ID, util.TracksToIDs(tracks)...)
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
