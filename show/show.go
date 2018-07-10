package show

import (
	"github.com/brianloveswords/spotify/auth"
	"github.com/brianloveswords/spotify/logger"
	"github.com/brianloveswords/spotify/util"
	"github.com/fatih/color"
)

var glog = logger.DefaultLogger

func Show() {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient())
	name := util.SongAttributionFromTrack(track)
	glog.Log("currently playing %s", color.CyanString(name))
}
func Artist() {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient())
	glog.Log(track.Artists[0].Name)
}
func ArtistID() {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient())
	glog.Log("%s", track.Artists[0].ID)
}
func ArtistURI() {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient())
	glog.Log("%s", track.Artists[0].URI)
}
func Track() {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient())
	glog.Log("%s", track.Name)
}
func TrackID() {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient())
	glog.Log("%s", track.ID)
}
func TrackURI() {
	track := util.MustGetCurrentlyPlaying(auth.SetupClient())
	glog.Log("%s", track.URI)
}
