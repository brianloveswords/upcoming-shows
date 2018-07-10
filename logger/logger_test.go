package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerVisual(t *testing.T) {
	logger := New()

	logger.Level = LevelDebug

	logger.Debug("kick")
	logger.Debug("out")
	logger.Debug("the")
	logger.Debug("jams")
	logger.Debug("motherfucker")

	var (
		level1 func()
		level2 func()
		level3 func()
		level4 func()
	)

	level1 = func() {
		defer logger.Enter("level1")()

		logger.Debug("in level1")
		logger.Debug("before level2")
		level2()
		logger.Debug("after level2")
	}
	level2 = func() {
		defer logger.Enter("level2")()
		logger.Debug("in level2")
		logger.Debug("before level3")
		level3()
		logger.Debug("after level3")
	}
	level3 = func() {
		defer logger.Enter("level3")()
		logger.Debug("in level3")
		logger.Debug("before level4")
		level4()
		logger.Debug("after level4")
	}
	level4 = func() {
		defer logger.Enter("level4")()
		logger.Debug("in level4")
		logger.Debug("done")
	}

	level1()

	logger.Level = LevelExtreme
	level1()

	// shouldn't see output
	logger.Level = LevelSilent
	logger.Debug("if you're seeing this something is wrong")
}

func TestLoggerLevels(t *testing.T) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	log := Logger{
		Stderr: &stderr,
		Stdout: &stdout,
	}

	// everything but extreme should log
	logemall := func() {
		log.Extreme("extreme")
		log.Debug("debug")
		log.Verbose("verbose")
		log.Log("normal")
	}

	clear := func() {
		stdout = bytes.Buffer{}
		stderr = bytes.Buffer{}
	}

	// cmdoutput should go to stdout, not stderr
	log.CmdOutput("output")
	assert.Equal(t, "output\n", stdout.String())

	log.Level = LevelExtreme
	logemall()
	assert.Equal(t, "<top>: extreme\n<top>: debug\nverbose\nnormal\n", stderr.String())
	clear()

	log.Level = LevelDebug
	logemall()
	assert.Equal(t, "<top>: debug\nverbose\nnormal\n", stderr.String())
	clear()

	log.Level = LevelVerbose
	logemall()
	assert.Equal(t, "verbose\nnormal\n", stderr.String())
	clear()

	log.Level = LevelNormal
	logemall()
	assert.Equal(t, "normal\n", stderr.String())
	clear()

	log.Level = LevelSilent
	logemall()
	assert.Equal(t, "", stderr.String())
	clear()
}
