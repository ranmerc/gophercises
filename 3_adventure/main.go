package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"
)

// Option is an arc option. Arc is the the arc to go to next and text is shown to user.
type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

// Arc is a story arc.
type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

// Arcs is a set of story arcs.
type Arcs map[string]Arc

// ArcHandler handles the http request for a particular arc.
type ArcHandler struct {
	Pattern  string
	Title    string
	Arc      Arc
	Template *template.Template
}

// ServeHTTP ...
func (a ArcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Template.Execute(w, a.Arc)
}

// fatal is used when an unrecoverable error is encountered.
func fatal(msg string, err error) {
	fmt.Print(msg, err)
	os.Exit(0)
}

func main() {
	file, err := os.ReadFile("gopher.json")
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

		arcHandlers = append(arcHandlers, ArcHandler{
			Template: tmpl,
			Pattern:  pattern,
			Title:    arc.Title,
			Arc:      arc,
		})
	}

	for _, handler := range arcHandlers {
		mux.Handle(handler.Pattern, handler)
	}

	fmt.Println("Serving on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
