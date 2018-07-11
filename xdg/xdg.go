package xdg

import (
	"os"
	"path"

	"github.com/spf13/afero"
)

type App struct {
	Home  string
	App   string
	AppFs afero.Fs
}

var DataDir = ".local/share"
var ConfigDir = ".config"
var CacheDir = ".cache"

func NewApp(name string) (*App, error) {
	app := App{
		Home:  os.Getenv("HOME"),
		App:   name,
		AppFs: afero.NewOsFs(),
	}
	if err := app.MakeDirs(); err != nil {
		return nil, err
	}
	return &app, nil
}

func (a *App) MakeDirs() (err error) {
	var p string
	p = path.Join(a.Home, a.App, DataDir)
	if err = a.AppFs.MkdirAll(p, 0700); err != nil {
		return err
	}
	p = path.Join(a.Home, a.App, CacheDir)
	if err = a.AppFs.MkdirAll(p, 0700); err != nil {
		return err
	}
	p = path.Join(a.Home, a.App, ConfigDir)
	if err = a.AppFs.MkdirAll(p, 0700); err != nil {
		return err
	}
	return nil
}

var flags = os.O_RDWR | os.O_CREATE | os.O_TRUNC

func (a *App) dataFile(name string) string {
	return path.Join(a.Home, DataDir, a.App, name)
}
func (a *App) DataCreate(name string) (afero.File, error) {
	return a.AppFs.OpenFile(a.dataFile(name), flags, 0700)
}
func (a *App) DataOpen(name string) (afero.File, error) {
	return a.AppFs.Open(a.dataFile(name))
}
func (a *App) DataRemove(name string) error {
	return a.AppFs.Remove(a.dataFile(name))
}

func (a *App) configFile(name string) string {
	return path.Join(a.Home, ConfigDir, a.App, name)
}
func (a *App) ConfigCreate(name string) (afero.File, error) {
	return a.AppFs.OpenFile(a.configFile(name), flags, 0700)
}
func (a *App) ConfigOpen(name string) (afero.File, error) {
	return a.AppFs.Open(a.configFile(name))
}
func (a *App) ConfigRemove(name string) error {
	return a.AppFs.Remove(a.configFile(name))
}

func (a *App) cacheFile(name string) string {
	return path.Join(a.Home, CacheDir, a.App, name)
}
func (a *App) CacheCreate(name string) (afero.File, error) {
	return a.AppFs.OpenFile(a.cacheFile(name), flags, 0700)
}
func (a *App) CacheOpen(name string) (afero.File, error) {
	return a.AppFs.Open(a.cacheFile(name))
}
func (a *App) CacheRemove(name string) error {
	return a.AppFs.Remove(a.cacheFile(name))
}
