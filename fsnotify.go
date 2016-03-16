package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func NewWatcher(file string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(file)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

var Usage = func() {
	fmt.Fprintf(os.Stdout, "Usage:  %s filename\n\n", os.Args[0])

	os.Exit(1)
}

// FileExists returns whether a file exists at a given filesystem path.
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, fmt.Errorf("not found: %s", path)
	}
	return false, err
}

var filename string

func init() {
	flag.Usage = Usage
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
	}

	filename = flag.Arg(0)

	if ret, err := FileExists(filename); !ret {
		log.Println(err)
		flag.Usage()
	}
}

func main() {
	filename := flag.Arg(0)
	NewWatcher(filename)
}
