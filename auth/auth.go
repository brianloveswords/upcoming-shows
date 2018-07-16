package auth

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/brianloveswords/spotify/logger"
	"github.com/brianloveswords/spotify/util"
	"github.com/lpabon/godbc"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

const redirectURL = "http://localhost:8888"

var permissions = []string{
	spotify.ScopeUserReadPrivate,
	spotify.ScopeUserLibraryRead,
	spotify.ScopePlaylistModifyPrivate,
	spotify.ScopePlaylistModifyPublic,
	spotify.ScopeUserLibraryModify,
	spotify.ScopeUserModifyPlaybackState,
	spotify.ScopeUserReadCurrentlyPlaying,
	spotify.ScopeUserReadRecentlyPlayed,
}

var glog = logger.DefaultLogger
var client *spotify.Client

func SetupClient() *spotify.Client {
	// TODO: if the token is older than a certain timeframe, force revalidation

	defer glog.Enter("auth.SetupClient")()

	// the redirect URL must be an exact match of a URL you've registered for your application
	// scopes determine which permissions the user is prompted to authorize
	var tok *oauth2.Token

	auth := spotify.NewAuthenticator(redirectURL, permissions...)

	id, secret := clientID, clientSecret
	auth.SetAuthInfo(id, secret)

	// see if we can just load a token straight up
	tok = loadToken()

	if tok != nil {
		expiry := tok.Expiry
		defer func() {
			if newtok, err := client.Token(); err == nil {
				if newtok.Expiry.After(expiry) {
					glog.Debug("saving new token: %v", tok.Expiry)
					saveToken(tok)
				}
			}
		}()

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

	util.OpenURL(url, false)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// use the same state string here that you used to generate the URL
		token, err := auth.Token(state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusNotFound)
			return
		}
		// save the token and create that shizz
		saveToken(token)
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

func randomState() string {
	b := make([]byte, 24)
	io.ReadFull(rand.Reader, b)
	return hex.EncodeToString(b)
}
