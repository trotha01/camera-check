package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLastAlertTimeZeroTime(t *testing.T) {
	n := Notifier{}
	lastTime, err := n.lastAlertTime()
	assert.Nil(t, err)
	t.Log(lastTime)
	assert.True(t, lastTime.IsZero())
}

func TestLastAlertTimeParseError(t *testing.T) {
	tmpfile := tmpFile(t, "not a real time")
	defer tmpFileCleanup(t, tmpfile)

	n := Notifier{
		lastAlertFile: tmpfile.Name(),
	}
	lastTime, err := n.lastAlertTime()
	assert.True(t, lastTime.IsZero())
	require.NotNil(t, err)
	_, ok := err.(*time.ParseError)
	assert.True(t, ok)
}

func TestLastAlertTime(t *testing.T) {
	now := time.Now().UTC()
	tmpfile := tmpFile(t, now.Format(timeFormat))
	defer tmpFileCleanup(t, tmpfile)
	n := Notifier{timeFormat: timeFormat,
		lastAlertFile: tmpfile.Name(),
	}
	lastTime, err := n.lastAlertTime()
	assert.False(t, lastTime.IsZero())
	require.Nil(t, err)
	assert.Equal(t, now.Format(timeFormat), lastTime.Format(timeFormat),
		"Got: %s Want: %s",
		lastTime.Format(timeFormat),
		now.Format(timeFormat))
}

func tmpFile(t *testing.T, content string) *os.File {
	tmpFile, err := ioutil.TempFile(".", "tmpfileprefix")
	if err != nil {
		t.Logf("Error creating tmp file: %s", err)
		t.FailNow()
	}

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Logf("Error writing to tmp file: %s", err)
		t.FailNow()
	}

	return tmpFile
}

func tmpFileCleanup(t *testing.T, tmpFile *os.File) {
	if err := tmpFile.Close(); err != nil {
		t.Logf("Error writing to tmp file: %s", err)
		t.FailNow()
	}

	os.Remove(tmpFile.Name())
}
