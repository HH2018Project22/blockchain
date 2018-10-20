package main

import (
	"encoding/json"
	"flag"
	"log"

	"github.com/HH2018Project22/bloodcoin/blockchain"
)

var (
	prescriptionCommand = flag.NewFlagSet("prescription", flag.ExitOnError)
	data                string
)

func init() {
	prescriptionCommand.StringVar(&data, "data", data, "Prescription data")
}

func doPrescription(args []string) {

	if err := prescriptionCommand.Parse(args); err != nil {
		panic(err)
	}

	bc := getBlockchain()

	log.Println("adding prescription")

	prescription := &blockchain.Prescription{}
	if err := json.Unmarshal([]byte(data), prescription); err != nil {
		panic(err)
	}

	prescriptionEvent := blockchain.NewPrescriptionEvent(prescription)
	if result := bc.AddEvent(prescriptionEvent); !result {
		panic("invalid prescription")
	}

	log.Println("saving blockchain")
	if err := bc.Save(blockchainPath); err != nil {
		panic(err)
	}

}
