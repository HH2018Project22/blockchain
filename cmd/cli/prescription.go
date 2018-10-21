package main

import (
	"encoding/json"
	"flag"
	"log"

	"github.com/HH2018Project22/bloodcoin/blockchain"
	"github.com/btcsuite/btcutil/base58"
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
	block, err := bc.AddEvent(prescriptionEvent)

	if err != nil {
		panic(err)
	}

	log.Println("Block:", base58.Encode(block.Hash))
	log.Println("saving blockchain")
	if err := bc.Save(blockchainPath); err != nil {
		panic(err)
	}

}
