package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var songkickDataFilename = "artist-songkick.data"

func loadSongkickData() (intmap map[string]int) {
	return loadIntMap(songkickDataFilename)
}
func saveSongkickData(skmap map[string]int) {
	saveIntMap(songkickDataFilename, skmap)
}

func getIDFromSongkickPage(artist string) int {
	resp, err := http.Get(songkickArtistURL(artist))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return parseSongkickPage(resp.Body, artist)
}

func parseSongkickPage(src io.Reader, artist string) int {
	artist = strings.ToLower(artist)

	doc, err := html.Parse(src)
	if err != nil {
		panic(err)
	}

	var artistNodes []*html.Node
	forEachNode(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "artist" {
					artistNodes = append(artistNodes, n)
				}
			}
		}
	}, nil)

	for _, artistNode := range artistNodes {
		link := findArtistLink(artistNode)
		name := artistNameFromLink(link)

		if strings.ToLower(name) == artist {
			return idFromLink(link)
		}
	}

	return SongkickNotFound
}
func artistNameFromLink(n *html.Node) (name string) {
	forEachNode(n, func(n *html.Node) {
		if n.Type == html.TextNode {
			name = n.Data
		}
	}, nil)
	return name
}

func idFromLink(n *html.Node) (id int) {
	for _, a := range n.Attr {
		if a.Key == "href" {
			return idFromHref(a.Val)
		}
	}
	return SongkickNotFound
}

func idFromHref(href string) int {
	base := path.Base(href)
	fields := strings.Split(base, "-")
	id, err := strconv.Atoi(fields[0])
	if err != nil {
		return SongkickNotFound
	}
	return id
}

func findArtistLink(n *html.Node) *html.Node {
	context := false
	var visit func(n, link *html.Node) *html.Node
	visit = func(n, link *html.Node) *html.Node {
		if link != nil {
			return link
		}

		// here's the fragment we're examining
		//
		// <li class="artist">
		//   <a href="/artists/7180534-gleemer" class="thumb">
		//     <img src="..." width="74" height="74" alt="" class="profile-pic artist">
		//   </a>
		//   <div class="subject">
		//     <span class="item-state-tag search-result">Artist</span>
		//     <p class="summary">
		// 	     <a href="/artists/7180534-gleemer"><strong>Gleemer</strong></a>
		// 	     ...
		//
		// we want to find that *second* link, so we first look for a
		// div that has the class "subject" to indicate we're in the
		// right context before we evaluate if we've found a link

		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "subject" {
					context = true
				}
			}
		}

		if context && n.Type == html.ElementNode && n.Data == "a" {
			return n
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			link = visit(c, link)
		}
		return link
	}
	return visit(n, nil)
}

// forEachNode calls the functions pre(x) and post(x) for each node x in
// the tree rooted at n. Both functions are optional. pre is called
// before the children are visited (preorder) and post is called after
// (postorder)
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}

var SongkickUnknown = 0
var SongkickNotFound = -1

var songkickBaseURL = "https://www.songkick.com/search?utf8=✓&query=%s&type=artists"

func songkickArtistURL(name string) string {
	location := fmt.Sprintf(songkickBaseURL, url.QueryEscape(name))
	fmt.Println(location)
	return location
}

func lookupSongkickIDs(artists []Artist) {
	skmap := loadSongkickData()

	var notfound []Artist
	var manual []Artist

	for _, artist := range artists {
		artist.SongkickID = skmap[artist.Name]

		if artist.SongkickID == SongkickUnknown {
			notfound = append(notfound, artist)
		}
	}

	fmt.Printf("lookup up songkick IDs for %d artists...\n", len(notfound))

	for _, artist := range notfound {
		artist.SongkickID = skmap[artist.Name]

		if artist.SongkickID != SongkickUnknown {
			continue
		}

		// try to look up automatically
		id := getIDFromSongkickPage(artist.Name)
		if id == SongkickNotFound {
			manual = append(manual, artist)
			fmt.Printf("results unclear for %s, skipping...\n", artist.Name)
			continue
		}
		fmt.Printf("Songkick ID for %s: %d\n", artist.Name, id)
		artist.SongkickID = id
		skmap[artist.Name] = id
		saveSongkickData(skmap)

	}

	fmt.Printf("manual identification needed for %d artists:\n", len(notfound))

	for _, artist := range manual {
		// read from stdin until we get a valid input
		openURL(songkickArtistURL(artist.Name))
		for {
			fmt.Printf("Enter Songkick ID for %s: ", artist.Name)
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			text = strings.Trim(text, "\n ")

			if text == "" {
				fmt.Printf("marking %s as not found\n\n", artist.Name)
				artist.SongkickID = SongkickNotFound
				skmap[artist.Name] = SongkickNotFound
				saveSongkickData(skmap)
				break
			}

			id, err := strconv.Atoi(text)
			if err != nil {
				fmt.Printf("invalid ID %s, must be an int\n\n", text)
				continue
			}

			artist.SongkickID = id
			skmap[artist.Name] = id
			saveSongkickData(skmap)
			break

		}
	}
}
