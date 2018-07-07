package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"sort"
	"time"

	"github.com/lpabon/godbc"
	"github.com/zmb3/spotify"
)

const redirectURL = "http://localhost:8888"

// these will be set at build time from env SPOTIFY_ID and SPOTIFY_SECRET
var clientID string
var clientSecret string

var tokenPath = path.Join(os.Getenv("HOME"), ".upcoming-shows-token")

// allows runtime overriding of build time auth info
func loadAuthInfo() (id, secret string) {
	id = mustGet(clientID, "SPOTIFY_ID")
	secret = mustGet(clientSecret, "SPOTIFY_SECRET")
	return id, secret
}

func openURL(url string) {
	cmd := exec.Command("open", url)
	if err := cmd.Run(); err != nil {
		// fall back to just printing it
		fmt.Printf("go here and authenticate\n: %s\n", url)
		return
	}
}

func setupClient() (client *spotify.Client) {
	// the redirect URL must be an exact match of a URL you've registered for your application
	// scopes determine which permissions the user is prompted to authorize
	auth := spotify.NewAuthenticator(redirectURL, spotify.ScopeUserReadPrivate, spotify.ScopeUserLibraryRead)

	id, secret := loadAuthInfo()
	auth.SetAuthInfo(id, secret)

	// see if we can just load a token straight up
	tok := loadToken(tokenPath)

	if tok != nil {
		c := auth.NewClient(tok)
		client = &c
		godbc.Ensure(client != nil, "failed to create client")
		return client
	}

	var s *http.Server

	state := randomState()

	// if you didn't store your ID and secret key in the specified environment variables,
	// you can set them manually here
	// auth.SetAuthInfo(clientID, secretKey)

	// get the user to this URL - how you do that is up to you
	// you should specify a unique state string to identify the session
	url := auth.AuthURL(state)

	openURL(url)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// use the same state string here that you used to generate the URL
		token, err := auth.Token(state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusNotFound)
			return
		}
		// save the token and create that shizz
		saveToken(token, tokenPath)
		c := auth.NewClient(token)
		client = &c

		w.WriteHeader(200)
		w.Write([]byte("<html><body>cool thx<script>window.close()</script>"))

		// we need to let the handler function complete in order for the
		// writes to be flushed, so we stick the s.Close() in a
		// goroutine that waits just a hot (milli)second before closing
		// the server.
		go func() {
			time.Sleep(1 * time.Millisecond)
			s.Close()
		}()
	})

	s = &http.Server{
		Addr:           ":8888",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}

	godbc.Ensure(client != nil, "failed to create client")
	return client
}

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
	tracks := getAllTracks(client)

	fmt.Fprintf(os.Stderr, "user: %s\n", user.ID)
	artists := processTracklist(tracks)
	printHistogram(artists)
}

var SongkickUnknown = 0
var SongkickNotFound = -1

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

var trackDataFilename = "saved-tracks.data"

func mustCreate(filename string) *os.File {
	f, err := os.Create(trackDataFilename)
	if err != nil {
		panic(err.Error())
	}
	return f
}

func loadTrackData() (tracks []spotify.SavedTrack) {
	f, err := os.Open(trackDataFilename)
	if err != nil {
		debugprint("couldn't open saved track data from file '%s'\n", trackDataFilename)
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

type trackData struct {
	Tracks []spotify.SavedTrack
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
		debugprint("loading tracks from disk\n")
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

		debugprint("got %s", page.Endpoint)

		tracks = append(tracks, page.Tracks...)

		if len(tracks) >= page.Total {
			break
		}
	}
	saveTrackData(tracks)
	return tracks
}
