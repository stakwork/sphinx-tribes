package routes

import (
	"net/http"
	"runtime/trace"
	"log"
//	"os"
	"bytes"

	"github.com/go-chi/chi"
)

func TestRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/internal-server-error", func(w http.ResponseWriter, r *http.Request) {
					// Enable tracing
	//f, err := os.Create("trace.out")
	var buf bytes.Buffer
	//if err != nil {
	//	log.Fatalf("Failed to create trace output file: %v", err)
	//}
	//defer f.Close()

	if err := trace.Start(&buf); err != nil {
		log.Fatalf("Failed to start trace: %v", err)
	}
	defer func() {
					trace.Stop()
					log.Println("Trace Data:")
	log.Println(buf.String())
	}()
		//panic("Forced internal server error")
	})

	return r
}
