package main

import (
	"github.com/fatih/color"
)

func showShow() {
	track := mustGetCurrentlyPlaying()
	name := songAttributionFromTrack(track)
	glog.Log("currently playing %s", color.CyanString(name))
}

func showArtist() {
	track := mustGetCurrentlyPlaying()
	glog.Log(track.Artists[0].Name)
}
func showArtistID() {
	track := mustGetCurrentlyPlaying()
	glog.Log("%s", track.Artists[0].ID)
}
func showArtistURI() {
	track := mustGetCurrentlyPlaying()
	glog.Log("%s", track.Artists[0].URI)
}
func showTrack() {
	track := mustGetCurrentlyPlaying()
	glog.Log("%s", track.Name)
}
func showTrackID() {
	track := mustGetCurrentlyPlaying()
	glog.Log("%s", track.ID)
}
func showTrackURI() {
	track := mustGetCurrentlyPlaying()
	glog.Log("%s", track.URI)
}
