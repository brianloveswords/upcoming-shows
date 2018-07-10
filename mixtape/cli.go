package mixtape

import (
	"fmt"
	"strconv"

	"github.com/brianloveswords/spotify/auth"
	"github.com/brianloveswords/spotify/cli"
	"github.com/brianloveswords/spotify/util"
	"github.com/zmb3/spotify"
)

var (
	paramArtist   *string
	paramTrackID  *spotify.ID
	paramLength   int
	defaultLength = 10
)

var CLI = cli.Command{
	Name: "mixtape",
	Help: "create a mixtape",
	Func: func() {
		if paramArtist == nil && paramTrackID == nil {
			glog.Fatal("must pass artist or track parameter")
		}
		if paramArtist != nil {
			if *paramArtist != "" {
				ByArtist(*paramArtist, paramLength)
				return
			}

			track := util.MustGetCurrentlyPlaying(auth.SetupClient())
			artistID := track.Artists[0].ID
			ByArtistID(artistID, paramLength)
			return
		}

		if paramTrackID != nil {
			if *paramTrackID == spotify.ID("") {
				ByCurrentTrack(paramLength)
				return
			}
			ByTrackID(*paramTrackID, paramLength)
			return
		}
	},

	Examples: []cli.Example{
		cli.Example{
			Args: []string{`artist`},
			Desc: "create a mixtape from artist currently playing",
		},
		cli.Example{
			Args: []string{`artist`, `length=20`},
			Desc: "create a mixtape of 20 tracks from artist currently playing",
		},
		cli.Example{
			Args: []string{`artist="The Sword"`},
			Desc: "create a mixtape of songs by The Sword",
		},
		cli.Example{
			Args: []string{`artist=bill`},
			Desc: "if the artist search is ambiguous, you will be prompted with selections",
		},
		cli.Example{
			Args: []string{`track`},
			Desc: "create mixtape of tracks recommended based on current playing track",
		},
		cli.Example{
			Negative: true,
			Args:     []string{`artist=Chavez`, `track`},
			Comment:  "THIS WON'T WORK!",
			Desc:     "note `artist` and `track` are mutually exclusive, don't include both",
		},
	},
	Params: []cli.Param{
		cli.Param{
			Name:  "artist",
			Alias: []string{"a"},
			Help:  "artist to base mixtape on. pass empty artist to use current playing",
			ParseFn: func(val string) error {
				if paramTrackID != nil {
					return fmt.Errorf("must pass track or artist, but not both")
				}
				paramArtist = &val
				return nil
			},
		},
		cli.Param{
			Name:  "track",
			Alias: []string{"t"},
			Help:  "track ID to base mixtape on. pass empty track to use current playing",
			ParseFn: func(val string) error {
				if paramArtist != nil {
					return fmt.Errorf("must pass track or artist, but not both")
				}
				trackID := spotify.ID(val)
				paramTrackID = &trackID
				return nil
			},
		},
		cli.Param{
			Name:     "length",
			Alias:    []string{"n"},
			Help:     fmt.Sprintf("number of tracks to include. defaults to %d", defaultLength),
			Implicit: true,
			ParseFn: func(val string) (err error) {
				if val == "" {
					paramLength = defaultLength
					return nil
				}
				paramLength, err = strconv.Atoi(val)
				if paramLength < 1 || paramLength > 50 {
					return fmt.Errorf("mixtape length must be between 1 and 50")
				}
				return err
			},
		},
	},
}
