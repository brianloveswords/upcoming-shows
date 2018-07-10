package show

import "github.com/brianloveswords/spotify/cli"

var CLI = cli.Command{
	Name: "show",
	Help: "show the current playing track. has subcommands",
	Func: Show,
	Commands: cli.Subcommands{
		&cli.Command{
			Name: "artist",
			Help: "show the artist of the current track",
			Func: Artist,
		},
		&cli.Command{
			Name: "artist-id",
			Help: "show the spotify artist ID of current track",
			Func: ArtistID,
		},
		&cli.Command{
			Name: "artist-uri",
			Help: "show the spotify URI for the artist of current track",
			Func: ArtistURI,
		},
		&cli.Command{
			Name: "track",
			Help: "show the title for the current track",
			Func: Track,
		},
		&cli.Command{
			Name: "track-id",
			Help: "show the track ID of the current track",
			Func: TrackID,
		},
		&cli.Command{
			Name: "track-uri",
			Help: "show the spotify URI for the current track",
			Func: TrackURI,
		},
	},
}
