package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/Anixy/go-vue-js/ui"
)

func main() {
	srv := &http.Server{
		Addr:        ":8888",
		Handler:     router(),
		IdleTimeout: time.Minute,
	}

	srv.ListenAndServe()
}

func router() http.Handler {
	mux := http.NewServeMux()

	// index
	mux.HandleFunc("/", indexHandler)

	// static files
	staticFS, _ := fs.Sub(ui.StaticFiles, "dist")
	httpFS := http.FileServer(http.FS(staticFS))
	mux.Handle("/static/", httpFS)

	// api
	mux.HandleFunc("/api/v1/greeting", greetingAPI)
	return mux
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api") {
		http.NotFound(w, r)
		return
	}

	if r.URL.Path == "/favicon.ico" {
		rawFile, _ := ui.StaticFiles.ReadFile("dist/favicon.ico")
		w.Write(rawFile)
		return
	}

	rawFile, _ := ui.StaticFiles.ReadFile("dist/index.html")
	w.Write(rawFile)
}

func greetingAPI(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, there!"))
}