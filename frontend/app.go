package frontend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gobuffalo/packr/v2"
)

var appBox = packr.New("app", "./app/build")

func PingRoute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("pong")
}

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

// FaviconRoute favicon.ico
func FaviconRoute(w http.ResponseWriter, r *http.Request) {
	favicon, err := appBox.Find("favicon.ico")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(favicon))
}

// StaticRoute - css and js
func StaticRoute(w http.ResponseWriter, r *http.Request) {
	path := r.URL.RequestURI()
	// if strings.HasPrefix(path, "/t/") {
	// 	path = path[2:]
	// }
	if strings.HasPrefix(path, "/static/css") {
		w.Header().Set("content-type", "text/css")
	}
	if strings.HasPrefix(path, "/static/js") {
		w.Header().Set("content-type", "application/javascript")
	}
	if strings.HasSuffix(path, ".svg") {
		w.Header().Set("content-type", "image/svg+xml")
	}
	file, err := appBox.Find(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(file))
}
