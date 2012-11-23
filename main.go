// gohttpd project main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	root := "."
	port := "80"
	flag.Parse()
	if flag.NArg() == 2 {
		root = flag.Arg(0)
		port = flag.Arg(1)
	} else {
		fmt.Println("Usage: gohttpd.exe root_dir port (defaut: gohttpd . 80)")
	}
	fmt.Printf("START gohttpd (DIR: %s  PORT: %s )\n", root, port)
	start(root, port)
}

func start(root, port string) {
	root = strings.Replace(root, "\\", "/", -1)
	root = strings.TrimRight(root, "/") + "/"
	http.Handle("/", http.FileServer(http.Dir(root)))
	s := &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
