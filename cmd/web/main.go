package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	dir  = flag.String("dir", "out/", "directory to serve")
	port = flag.String("port", ":81", "port to serve from")
)

func main() {
	flag.Parse()
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(*dir))))
	log.Printf("static server started on port %s\n", *port)
	log.Fatalln(http.ListenAndServe(*port, nil))
}
