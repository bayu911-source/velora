// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Hello is a simple hello, world demonstration web server.
//
// It serves version information on /version and answers
// any other request like /name by saying "Hello, name!".
//
// See golang.org/x/example/outyet for a more sophisticated server.
package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: helloserver [options]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var (
	greeting = flag.String("g", "Hello", "Greet with `greeting`")
	addr     = flag.String("addr", "localhost:8080", "address to serve")
)

func helloHandler(greeting string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/")
		if name == "" {
			name = "world"
		}
		fmt.Fprintf(w, "%s, %s!", greeting, html.EscapeString(name))
	}
}

func versionHandler(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, version)
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()

	version := "unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		version = info.Main.Version
	}

	http.HandleFunc("/version", versionHandler(version))
	http.HandleFunc("/", helloHandler(*greeting))

	log.Printf("serving on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
