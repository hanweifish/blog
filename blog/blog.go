// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command blog is a web server for the Go blog that can run on App Engine or
// as a stand-alone HTTP server.
package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/tools/blog"
	"golang.org/x/tools/godoc/static"

	_ "golang.org/x/tools/playground"
)

const hostname = "blog.golang.org" // default hostname for blog server

var (
	httpAddr = flag.String("http", "0.0.0.0:8080", "HTTP listen address")
	barContentPath = flag.String("bar", "bar/", "path to bar content files")
	contentPath = flag.String("content", "content/", "path to content files")
	templatePath = flag.String("template", "template/", "path to template files")
	staticPath = flag.String("static", "static/", "path to static files")

	config = blog.Config{
		Hostname:     hostname,
		BaseURL:      "//" + hostname,
		GodocURL:     "//golang.org",
		HomeArticles: 5, // articles to display on the home page
		FeedArticles: 10, // articles to include in Atom and JSON feeds
		PlayEnabled:  true,
		FeedTitle:    "Weekend Chinese Tech Meetup",
	}
)

func init() {
	// Redirect "/blog/" to "/", because the menu bar link is to "/blog/"
	// but we're serving from the root.
	redirect := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	http.HandleFunc("/blog", redirect)
	http.HandleFunc("/blog/", redirect)

	http.Handle("/lib/godoc/", http.StripPrefix("/lib/godoc/", http.HandlerFunc(staticHandler)))
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path
	b, ok := static.Files[name]
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, name, time.Time{}, strings.NewReader(b))
}

func serveBlog(prefix string, contentPath string) {
	blogConfig := blog.Config{
		Hostname: hostname,
		BaseURL: "//" + hostname,
		HomeArticles: 5,
		FeedArticles: 10,
		PlayEnabled: true,
		FeedTitle: "unimportant",
		ContentPath: contentPath,
		TemplatePath: *templatePath,
	}

	s2, err := blog.NewServer(blogConfig)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle(prefix + "/", http.StripPrefix(prefix, s2))
}

func main() {
	flag.Parse()

	serveBlog("/bar", *barContentPath)
	serveBlog("", *contentPath)

	log.Printf("xfguo: hello world\n")

	fs := http.FileServer(http.Dir(*staticPath))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
