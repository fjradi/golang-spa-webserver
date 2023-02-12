package main

import (
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	http.HandleFunc("/", handleSPA)
	http.ListenAndServe(":80", nil)
}

func handleSPA(w http.ResponseWriter, r *http.Request) {
	rootPath, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	buildPath := filepath.Join("production")
	path := filepath.Join(rootPath, buildPath, r.URL.Path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(buildPath, "index.html"))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(buildPath)).ServeHTTP(w, r)
}
