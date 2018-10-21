package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/HH2018Project22/bloodcoin/blockchain"
	"github.com/btcsuite/btcutil/base58"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var (
	dbPath       string = "./bloodcoin-server.db"
	bc           *blockchain.Blockchain
	peerEndpoint = ""
)

type PrescriptionHash struct {
	Hash string `json:"hash"`
}

func init() {
	flag.StringVar(&dbPath, "db", dbPath, "Database file path")
}

func main() {

	flag.Parse()

	var err error
	var syncHook blockchain.BlockHookFunc
	if peerEndpoint != "" {
		syncHook = blockchain.CreateBlockSyncHook(peerEndpoint)
	}

	bc, err = blockchain.LoadBlockchain(dbPath, syncHook)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Route("/prescriptions", func(r chi.Router) {
		r.Get("/", ListPrescriptions)
		r.Post("/new", CreatePrescription)
		r.Route("/{prescriptionHash}", func(r chi.Router) {
			r.Get("/", ReadPrescription)
			r.Get("/notifications", ReadPrescriptionNotifications)
		})
	})

	r.Route("/notifications", func(r chi.Router) {
		r.Post("/new", CreateNotification)
	})

	r.Route("/blocks", func(r chi.Router) {
		r.Get("/", ListBlocks)
		r.Post("/new", CreateBlock)
	})

	staticHandler := http.FileServer(http.Dir("web"))
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		staticHandler.ServeHTTP(w, r)
	})

	http.ListenAndServe(":3000", r)
}

func ListBlocks(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	writeJSON(w, bc.Blocks())
}

func ListPrescriptions(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	prescriptions := bc.ListPrescriptions()
	writeJSON(w, prescriptions)
}

func ReadPrescription(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	prescriptionHash := chi.URLParam(r, "prescriptionHash")

	b58Hash := base58.Decode(prescriptionHash)
	block := bc.FindPrescriptionBlock(b58Hash)

	if block == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	prescriptionEvent := block.Event.(*blockchain.PrescriptionEvent)

	writeJSON(w, prescriptionEvent.Prescription)
}

func ReadPrescriptionNotifications(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	prescriptionHash := chi.URLParam(r, "prescriptionHash")

	b58Hash := base58.Decode(prescriptionHash)

	events := bc.FindPrescriptionNotificationEvents(b58Hash)

	s := make([]*blockchain.NotificationEvent, 0)

	for _, event := range events {
		s = append(s, event)
	}

	writeJSON(w, s)
}

// CreatePrescription persists the posted Prescription and returns it
// back to the client as an acknowledgement.
func CreatePrescription(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	decoder := json.NewDecoder(r.Body)

	var t = map[string]interface{}{}
	err := decoder.Decode(&t)

	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("result.csv", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	for k, v := range t {
		line := fmt.Sprintf("%s, %s\n", k, v)
		file.WriteString(line)
	}

	w.WriteHeader(http.StatusNoContent)

}

// CreateNotification adds a notification for a given Prescription in our blockchain.
func CreateNotification(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	notification := &blockchain.NotificationEvent{}
	parseJSON(r, notification)

	event := blockchain.NewNotificationEvent(
		notification.PrescriptionHash,
		notification.NotificationType,
		notification.Operator,
	)

	if _, err := bc.AddEvent(event); err != nil {
		panic(err)
	}

	writeJSON(w, event)

}

func CreateBlock(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	block := &blockchain.Block{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(block); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := bc.AddBlock(block); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func parseJSON(r *http.Request, data interface{}) {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(data); err != nil {
		panic(err)
	}
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		panic(err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
