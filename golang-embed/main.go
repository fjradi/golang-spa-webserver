package main

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed production
var buildFS embed.FS

func main() {
	http.HandleFunc("/", handleSPA)
	http.ListenAndServe(":80", nil)
}

func handleSPA(w http.ResponseWriter, r *http.Request) {
	buildPath := filepath.Join("production")
	path, err := buildFS.Open(filepath.Join(buildPath, r.URL.Path))
	if os.IsNotExist(err) {
		index, err := buildFS.ReadFile(filepath.Join(buildPath, "index.html"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write(index)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer path.Close()
	staticFS, err := fs.Sub(buildFS, buildPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.FS(staticFS)).ServeHTTP(w, r)
}
