package main

import (
	"context"
	"errors"
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

type Prescription struct {
	ID string `json:"id"`
}

type PrescriptionRequest struct {
	Prescription *Prescription `json:"prescription"`
}

type PrescriptionResponse struct {
	*Prescription

	// We add an additional field to the response here.. such as this
	// elapsed computed property
	Elapsed int64 `json:"elapsed"`
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

	r.Route("/prescriptions", func(r chi.Router) {
		//r.With(paginate).Get("/", ListPrescriptions)
		r.Post("/", CreatePrescription) // POST /prescriptions
		//r.Get("/search", SearchPrescriptions) // GET /prescriptions/search

		r.Route("/{prescriptionID}", func(r chi.Router) {
			r.Use(PrescriptionCtx)         // Load the *Prescription on the request context
			r.Get("/", GetPrescription)    // GET /prescriptions/123
			r.Put("/", UpdatePrescription) // PUT /prescriptions/123
		})

		// GET /prescriptions/whats-up
		//r.With(PrescriptionCtx).Get("/{prescriptionSlug:[a-z-]+}", GetPrescription)
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

// GetPrescription returns the specific Prescription. You'll notice it just
// fetches the Prescription right off the context, as its understood that
// if we made it this far, the Prescription must be on the context. In case
// its not due to a bug, then it will panic, and our Recoverer will save us.
func GetPrescription(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the prescription
	// context because this handler is a child of the PrescriptionCtx
	// middleware. The worst case, the recoverer middleware will save us.
	prescription := r.Context().Value("prescription").(*Prescription)

	if err := render.Render(w, r, NewPrescriptionResponse(prescription)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// CreatePrescription persists the posted Prescription and returns it
// back to the client as an acknowledgement.
func CreatePrescription(w http.ResponseWriter, r *http.Request) {
	data := &PrescriptionRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	prescription := data.Prescription
	blockchainNewPrescription(prescription)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewPrescriptionResponse(prescription))
}

// UpdatePrescription updates an existing Prescription in our persistent store.
func UpdatePrescription(w http.ResponseWriter, r *http.Request) {
	prescription := r.Context().Value("prescription").(*Prescription)

	data := &PrescriptionRequest{Prescription: prescription}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	prescription = data.Prescription
	blockchainUpdatePrescription(prescription)

	render.Render(w, r, NewPrescriptionResponse(prescription))
}

func PrescriptionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var prescription *Prescription
		var err error

		if prescriptionID := chi.URLParam(r, "prescriptionID"); prescriptionID != "" {
			prescription = blockchainGetPrescription(prescriptionID)
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

func blockchainGetPrescription(id string) (prescription *Prescription) {
	//TODO fetch Prescription from blockchain
	fmt.Printf("blockchainGetPrescription")
	fmt.Println()

	pres := NewPrescription()

	return pres
}

func NewPrescription() *Prescription {
	return &Prescription{
		ID: "1",
	}
}

func blockchainNewPrescription(prescription *Prescription) {
	//TODO add Prescription to blockchain
	fmt.Printf("blockchainNewPrescription")
	fmt.Println()
}

func blockchainUpdatePrescription(prescription *Prescription) {
	//TODO add Prescription event to blockchain
	fmt.Printf("blockchainUpdatePrescription")
	fmt.Println()
}

func (p *PrescriptionRequest) Bind(r *http.Request) error {
	// a.Prescription is nil if no Prescription fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if p.Prescription == nil {
		return errors.New("missing required Prescription fields.")
	}

	// just a post-process after a decode..
	return nil
}

func NewPrescriptionResponse(prescription *Prescription) *PrescriptionResponse {
	resp := &PrescriptionResponse{Prescription: prescription}

	return resp
}

func (rd *PrescriptionResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10
	return nil
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
