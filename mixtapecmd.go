package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/zmb3/spotify"
)

func mixtapeCreate() {
	if paramMixtapeArtist == nil && paramMixtapeTrackID == nil {
		glog.Fatal("must pass artist or track parameter")
	}
	if paramMixtapeArtist != nil {
		if *paramMixtapeArtist != "" {
			mixtapeByArtist(*paramMixtapeArtist, paramMixtapeLength)
			return
		}

		track := mustGetCurrentlyPlaying()
		artistID := track.Artists[0].ID
		mixtapeByArtistID(artistID, paramMixtapeLength)
		return
	}

	if paramMixtapeTrackID != nil {
		glog.Fatal("need to implement by track name")
	}
}

func mixtapeByArtist(artist string, length int) {
	var artistID spotify.ID
	client := setupClient()
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
		mixtapeByArtistID(artistID, paramMixtapeLength)
		return
	}

	for _, found := range artists {
		if strings.ToLower(found.Name) == normalizedArtist {
			artistID = found.ID
			mixtapeByArtistID(artistID, paramMixtapeLength)
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

	mixtapeByArtistID(pick.ID, paramMixtapeLength)
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

func mixtapeByArtistID(artistID spotify.ID, length int) {
	defer glog.Enter("mixtapeByArtistID")()
	client := setupClient()

	artist, err := client.GetArtist(artistID)
	if err != nil {
		glog.Fatal("couldn't look up artist with ID %s: %s", artistID, err)
	}

	alltracks, err := getAllTracksByArtist(client, artistID)
	if err != nil {
		glog.Fatal("could not get tracks from artist with ID %s: %s", artistID, err)
	}

	tracks := randomTracks(alltracks, length)
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
