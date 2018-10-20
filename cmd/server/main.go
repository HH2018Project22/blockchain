package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/docgen"
	"github.com/go-chi/render"
	"github.com/oxtoacart/bpool"
)

// Article data model. I suggest looking at https://upper.io for an easy
// and powerful data persistence adapter.
type Prescription struct {
	ID     string `json:"id"`
	UserID int64  `json:"user_id"` // the author
	Title  string `json:"title"`
	Slug   string `json:"slug"`
}

type TemplateConfig struct {
	TemplateLayoutPath  string
	TemplateIncludePath string
}

var routes = flag.Bool("routes", false, "Generate router documentation")

var templates map[string]*template.Template
var bufpool *bpool.BufferPool

var mainTmpl = `{{define "main" }} {{ template "base" . }} {{ end }}`
var templateConfig TemplateConfig

func loadConfiguration() {
	templateConfig.TemplateLayoutPath = "cmd/server/templates/layouts/"
	templateConfig.TemplateIncludePath = "cmd/server/templates/"
}

func main() {
	loadConfiguration()
	loadTemplates()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/prescription/new", prescriptionNew)

	r.Post("/prescription/create", prescriptionCreate)

	r.Route("/prescriptions", func(r chi.Router) {
		//r.With(paginate).Get("/", ListPrescriptions)
		r.Post("/", CreatePrescription) // POST /prescriptions
		//r.Get("/search", SearchPrescriptions) // GET /prescriptions/search

		r.Route("/{prescriptionID}", func(r chi.Router) {
			r.Use(PrescriptionCtx)         // Load the *Prescription on the request context
			r.Get("/", GetPrescription)    // GET /prescriptions/123
			r.Put("/", UpdatePrescription) // PUT /prescriptions/123
		})

		// GET /articles/whats-up
		//r.With(ArticleCtx).Get("/{articleSlug:[a-z-]+}", GetArticle)
	})

	// Passing -routes to the program will generate docs for the above
	// router definition. See the `routes.json` file in this folder for
	// the output.
	if *routes {
		// fmt.Println(docgen.JSONRoutesDoc(r))
		fmt.Println(docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			ProjectPath: "github.com/go-chi/chi",
			Intro:       "Welcome to the chi/_examples/rest generated docs.",
		}))
		return
	}

	http.ListenAndServe(":3000", r)
}

func loadTemplates() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	layoutFiles, err := filepath.Glob(templateConfig.TemplateLayoutPath + "*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	includeFiles, err := filepath.Glob(templateConfig.TemplateIncludePath + "*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	mainTemplate := template.New("main")

	mainTemplate, err = mainTemplate.Parse(mainTmpl)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range includeFiles {
		fileName := filepath.Base(file)
		files := append(layoutFiles, file)
		templates[fileName], err = mainTemplate.Clone()
		if err != nil {
			log.Fatal(err)
		}
		templates[fileName] = template.Must(templates[fileName].ParseFiles(files...))
	}

	log.Println("templates loading successful")

	bufpool = bpool.NewBufferPool(64)
	log.Println("buffer allocation successful")
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, fmt.Sprintf("The template %s does not exist.", name),
			http.StatusInternalServerError)
	}

	buf := bufpool.Get()
	defer bufpool.Put(buf)

	err := tmpl.Execute(buf, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

// paginate is a stub, but very possible to implement middleware logic
// to handle the request params for handling a paginated request.
func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}

func PrescriptionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var prescription *Prescription
		var err error

		if prescriptionID := chi.URLParam(r, "prescriptionID"); prescriptionID != "" {
			prescription, err = blockchainGetPrescription(prescriptionID)
		} else {
			render.Render(w, r, ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "prescription", prescription)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//--
// Error response payloads & renderers
//--

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
