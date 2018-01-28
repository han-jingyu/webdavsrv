package main

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/net/webdav"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {

	log.Println("== Simple Webdav Server ==")

	rootDir := os.Getenv("WEBDAV_ROOT_DIR")
	listenAddress := os.Getenv("WEBDAV_LISTEN_ADDRESS")

	if len(rootDir) == 0 {
		rootDir = "/webdav"
	}

	if len(listenAddress) == 0 {
		listenAddress = "localhost:8080"
	}

	kingpin.CommandLine.Help = "WEBDAV_ROOT_DIR, WEBDAV_LISTEN_ADDRESS are read from environment"
	kingpin.CommandLine.Help += ", though command line flags have priority"
	kingpin.Flag("rootdir", "directory to serve").Short('d').Default(rootDir).StringVar(&rootDir)
	kingpin.Flag("listen", "listen on address:port").Short('l').Default(listenAddress).StringVar(&listenAddress)
	kingpin.CommandLine.HelpFlag.Hidden()
	kingpin.Parse()

	log.Println("Config parameters:")
	log.Printf(" listenAddress: %s\n", listenAddress)
	log.Printf(" rootDir      : %s\n", rootDir)

	_, err := os.Stat(rootDir)
	if err != nil {
		log.Printf("Cannot use '%s' as webdav root: %v", err)
		os.Exit(1)
	}

	srv := &webdav.Handler{
		FileSystem: webdav.Dir(rootDir),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Printf("WEBDAV: %6s %s, ERROR: %v", r.Method, r.URL.Path, err)
			} else {
				log.Printf("WEBDAV: %6s %s", r.Method, r.URL.Path)
			}
		},
	}
	http.Handle("/", srv)
	log.Printf("Start serving on %s\n", listenAddress)
	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		log.Fatalf("Error with WebDAV server: %v", err)
	}
}
