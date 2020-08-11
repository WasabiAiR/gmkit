package notification

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matcornic/hermes/v2"
	"github.com/pkg/errors"
)

// DefaultPreviewer is the default previewer.
var DefaultPreviewer = &Previewer{
	Prefix: "/emailpreviews",
}

// Previewer is an HTTP handler that renders previews of email templates.
type Previewer struct {
	router    *mux.Router
	Data      map[string]hermes.Email
	Prefix    string
	URLPrefix string
	Renderer  Renderer
}

// Register registers a an email template with the given name.
func (p *Previewer) Register(name string, email hermes.Email) {
	if p.Data == nil {
		p.Data = make(map[string]hermes.Email)
	}

	p.Data[name] = email
}

func (p *Previewer) buildRoutes() {
	p.router = mux.NewRouter()
	p.router.HandleFunc("/{template}/{format}", p.emailPreview).Methods(http.MethodGet)
	p.router.HandleFunc("/", p.emailPreviewList).Methods(http.MethodGet)
}

// ServeHTTP builds the routes if necessary and serves up the request.
func (p *Previewer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" {
		r.URL.Path = "/"
	}
	if p.router == nil {
		p.buildRoutes()
	}

	p.router.ServeHTTP(w, r)
}

func (p *Previewer) emailPreviewList(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").Parse(tmplList)
	if err != nil {
		log.Println(errors.Wrap(err, "parsing template"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, p); err != nil {
		log.Println(errors.Wrap(err, "executing template"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func (p *Previewer) emailPreview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateName := vars["template"]
	format := vars["format"]

	email, ok := p.Data[templateName]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	txt, html, err := p.Renderer.Render(email)
	if err != nil {
		log.Println(errors.Wrap(err, "rendering email template"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch format {
	case "html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	default:
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(txt))
	}
}

const tmplList = `
<html>
<body>
<h1>Select a template to preview</h1>
<ul>
{{ $prefix := .Prefix }}
{{ $urlprefix := .URLPrefix }}
{{range $key, $v := .Data}}
<li><a href="{{$urlprefix}}{{$prefix}}/{{$key}}/html">{{$key}} - html</a></li>
<li><a href="{{$urlprefix}}{{$prefix}}/{{$key}}/txt">{{$key}} - txt</a></li>
{{end}}
</ul>
</body>
</html>`
