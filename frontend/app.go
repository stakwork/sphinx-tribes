package frontend

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gobuffalo/packr/v2"
)

var appBox = packr.New("app", "./app/build")

// IndexRoute index.html
func IndexRoute(w http.ResponseWriter, r *http.Request) {
	indexHTML, err := appBox.Find("index.html")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(indexHTML))
}

// ManifestRoute manifest.json
func ManifestRoute(w http.ResponseWriter, r *http.Request) {
	manifest, err := appBox.Find("manifest.json")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(manifest))
}

// StaticRoute - css and js
func StaticRoute(w http.ResponseWriter, r *http.Request) {
	path := r.URL.RequestURI()
	if strings.HasPrefix(path, "/static/css") {
		w.Header().Set("content-type", "text/css")
	}
	if strings.HasPrefix(path, "/static/js") {
		w.Header().Set("content-type", "application/javascript")
	}
	file, err := appBox.Find(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(file))
}
