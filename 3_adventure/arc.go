package main

import (
	"net/http"
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
	Arc      Arc
	Template *template.Template
}

// HandlerOptions is the type for the ArcHandler functional option.
type HandlerOption func(a *ArcHandler)

// WithTemplates is an functional option to add template to the ArcHandler.
func WithTemplate(t *template.Template) HandlerOption {
	return func(a *ArcHandler) {
		a.Template = t
	}
}

// NewArcHandler creates a new ArcHandler.
func NewArcHandler(pattern string, arc Arc, opts ...HandlerOption) ArcHandler {
	tmpl := template.Must(template.New("name").Parse("I am blank slate!"))

	handler := ArcHandler{
		Pattern:  pattern,
		Arc:      arc,
		Template: tmpl,
	}

	for _, opt := range opts {
		opt(&handler)
	}

	return handler
}

// ServeHTTP ...
func (a ArcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := a.Template.Execute(w, a.Arc); err != nil {
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
	}
}
