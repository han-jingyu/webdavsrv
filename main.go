package main

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/net/webdav"
)

func main() {

	log.Println("== Simple Webdav Server ==")

	rootDir := os.Getenv("WEBDAV_ROOT_DIR")
	listenAddress := os.Getenv("WEBDAV_LISTEN_ADDRESS")

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
