package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
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

func setupClient() (client *spotify.Client) {
	// the redirect URL must be an exact match of a URL you've registered for your application
	// scopes determine which permissions the user is prompted to authorize
	auth := spotify.NewAuthenticator(redirectURL, spotify.ScopeUserReadPrivate)

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

	fmt.Printf("go here and authenticate\n: %s\n", url)
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
		w.Write([]byte("cool thx bro"))

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

func main() {
	client := setupClient()
	if tok, err := client.Token(); err == nil {
		debugprint("token expires: %v", tok.Expiry)
		saveToken(tok, tokenPath)
	}
	fmt.Println(client.CurrentUser())
}
