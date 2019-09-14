package main

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

func FetchTokenTime(fname string) time.Time {
	var tm time.Time
	var x *exif.Exif
	f, err := os.Open(fname)
	if err != nil {
		return time.Now()
	}

	if isPic(fname) {
		x, err = exif.Decode(f)
		if err != nil {
			goto UseFileTime
		}
		tm, _ = x.DateTime()
		return tm
	} else if isMov(fname) {
		// TODO
		return time.Now()
	}

UseFileTime:
	fi, err := f.Stat()
	if err != nil {
		return time.Now()
	}

	tm = fi.ModTime()
	return tm
}

func isPic(fname string) bool {
	switch strings.ToLower(path.Ext(fname)) {
	case ".jpeg":
		fallthrough
	case ".jpg":
		fallthrough
	case ".png":
		return true
	}

	return false
}

func isMov(fname string) bool {
	switch strings.ToLower(path.Ext(fname)) {
	case ".mp4":
		return true
	}

	return false
}
