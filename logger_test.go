package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerVisual(t *testing.T) {
	logger := NewLogger()

	logger.Level = LoggerLevelDebug

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

	logger.Level = LoggerLevelExtreme
	level1()

	// shouldn't see output
	logger.Level = LoggerLevelSilent
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

	log.Level = LoggerLevelExtreme
	logemall()
	assert.Equal(t, "<top>: extreme\n<top>: debug\nverbose\nnormal\n", stderr.String())
	clear()

	log.Level = LoggerLevelDebug
	logemall()
	assert.Equal(t, "<top>: debug\nverbose\nnormal\n", stderr.String())
	clear()

	log.Level = LoggerLevelVerbose
	logemall()
	assert.Equal(t, "verbose\nnormal\n", stderr.String())
	clear()

	log.Level = LoggerLevelNormal
	logemall()
	assert.Equal(t, "normal\n", stderr.String())
	clear()

	log.Level = LoggerLevelSilent
	logemall()
	assert.Equal(t, "", stderr.String())
	clear()
}
