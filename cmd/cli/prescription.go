package main

import (
	"flag"
	"log"
	"time"

	"github.com/HH2018Project22/blockchain"
)

var (
	prescriptionCommand = flag.NewFlagSet("prescription", flag.ExitOnError)
	firstName           string
	lastName            string
	useName             string
	birthDate           time.Time
)

func doPrescription() {

	bc := getBlockchain()

	log.Println("adding prescription")

	prescription := blockchain.NewPrescriptionEvent(&blockchain.Patient{
		FirstName: "John",
		LastName:  "Doe",
		UseName:   "Doe",
	})

	bc.AddBlock(prescription)

	if err := bc.Save(blockchainPath); err != nil {
		log.Fatal(err)
	}

}
