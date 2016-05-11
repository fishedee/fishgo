package modules

import (
	"github.com/howeyc/fsnotify"
	"os"
	"path"
	"strings"
	"time"
)

type WatchCallback func(string)

var (
	watchFile     = map[string]bool{}
	watchCallback WatchCallback
	eventTime     = map[string]time.Time{}
)

func checkTMPFile(name string) bool {
	if strings.HasSuffix(strings.ToLower(name), ".tmp") {
		return true
	}
	return false
}

func checkIfWatchExt(name string) bool {
	if strings.HasSuffix(name, ".go") {
		return true
	}
	return false
}

func getFileModTime(path string) (time.Time, error) {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {
		return time.Time{}, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return time.Time{}, err
	}

	return fi.ModTime(), nil
}

func Watch(paths []string, callback WatchCallback) error {
	for _, path := range paths {
		watchFile[path] = true
	}
	watchCallback = callback
	return nil
}

func RunWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	for path, _ := range watchFile {
		Log.Debug("Watch Directory %v", path)
		err = watcher.Watch(path)
		if err != nil {
			panic(err)
			return err
		}
	}
	changeDirectory := ""
	newTimer := time.NewTimer(time.Second)
	watchCallback(".")
	for {
		select {
		case _ = <-newTimer.C:
			if changeDirectory != "" {
				watchCallback(changeDirectory)
				changeDirectory = ""
			}
			newTimer.Reset(time.Second)

		case e := <-watcher.Event:
			fileName := e.Name
			if checkTMPFile(fileName) {
				continue
			}

			if checkIfWatchExt(fileName) == false {
				continue
			}

			mt, err := getFileModTime(fileName)
			if err != nil {
				Log.Error("getFileModTime %v Fail", fileName)
				continue
			}
			if t := eventTime[fileName]; mt.Equal(t) {
				continue
			}
			eventTime[fileName] = mt

			Log.Debug("File Change %v", fileName)
			changeDirectory = path.Dir(fileName)
			newTimer.Reset(time.Second)

		case err := <-watcher.Error:
			Log.Error("%v", err.Error())
		}
	}
	return nil
}
