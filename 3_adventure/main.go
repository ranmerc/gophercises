package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"text/template"
)

// fatal is used when an unrecoverable error is encountered.
func fatal(msg string, err error) {
	fmt.Print(msg, err)
	os.Exit(0)
}

var jsonFile string
var port int

func init() {
	flag.IntVar(&port, "port", 8080, "the port to run the web server on")
	flag.StringVar(&jsonFile, "file", "gopher.json", "the JSON file name with the story in it")
}

func main() {
	flag.Parse()

	file, err := os.ReadFile(jsonFile)
	if err != nil {
		fatal("Failed to read json file", err)
	}

	arcs := Arcs{}
	if err = json.Unmarshal(file, &arcs); err != nil {
		fatal("Invalid json file", err)
	}

	tmpl, err := template.ParseFiles("arc.html")
	if err != nil {
		fatal("Invalid template file", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/intro", http.StatusPermanentRedirect)
	})

	arcHandlers := []ArcHandler{}
	for name, arc := range arcs {
		pattern := fmt.Sprintf("/%s", name)

		arcHandlers = append(arcHandlers, NewArcHandler(pattern, arc, WithTemplate(tmpl)))
	}

	for _, handler := range arcHandlers {
		mux.Handle(handler.Pattern, handler)
	}

	fmt.Printf("Serving on http://localhost:%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
