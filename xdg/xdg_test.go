package xdg

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestXDGApp(t *testing.T) {
	var err error
	testfs := afero.NewMemMapFs()
	basedir := App{
		Home:  "/",
		App:   "test-app",
		AppFs: testfs,
	}

	err = basedir.MakeDirs()
	assert.NoError(t, err)

	type testcase struct {
		Name   string
		Dir    string
		Create func(string) (afero.File, error)
		Open   func(string) (afero.File, error)
		Remove func(string) error
	}

	for _, tc := range []testcase{{
		Name:   "data",
		Dir:    "/.local/share/test-app",
		Create: basedir.DataCreate,
		Open:   basedir.DataOpen,
		Remove: basedir.DataRemove,
	}, {
		Name:   "config",
		Dir:    "/.config/test-app",
		Create: basedir.ConfigCreate,
		Open:   basedir.ConfigOpen,
		Remove: basedir.ConfigRemove,
	}, {
		Name:   "cache",
		Dir:    "/.cache/test-app",
		Create: basedir.CacheCreate,
		Open:   basedir.CacheOpen,
		Remove: basedir.CacheRemove,
	}} {
		t.Run(tc.Name, func(t *testing.T) {
			var f afero.File
			name := "test-data"
			f, err = tc.Create(name)
			assert.NoError(t, err)
			f.WriteString("lol")
			f.Close()

			assert.Equal(t, path.Join(tc.Dir, name), f.Name())

			f, err = tc.Open(name)
			assert.NoError(t, err)
			b, err := ioutil.ReadAll(f)
			assert.NoError(t, err)

			assert.Equal(t, "lol", string(b))

			err = tc.Remove(name)
			assert.NoError(t, err)

			_, err = testfs.Stat(path.Join(tc.Dir, name))
			assert.Error(t, err)
		})
	}

}
