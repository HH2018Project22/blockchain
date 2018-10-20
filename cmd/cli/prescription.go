package main

import (
	"flag"
	"log"
	"time"

	"github.com/HH2018Project22/bloodcoin/blockchain"
)

var (
	prescriptionCommand = flag.NewFlagSet("prescription", flag.ExitOnError)
	firstName           string
	lastName            string
	useName             string
	birthDate           string
)

func init() {
	prescriptionCommand.StringVar(&firstName, "first-name", firstName, "First name")
	prescriptionCommand.StringVar(&lastName, "last-name", lastName, "Last name")
	prescriptionCommand.StringVar(&useName, "use-name", useName, "Use name")
	prescriptionCommand.StringVar(&birthDate, "birth-date", birthDate, "Birth name")
}

func doPrescription(args []string) {

	if err := prescriptionCommand.Parse(args); err != nil {
		panic(err)
	}

	bc := getBlockchain()

	log.Println("adding prescription")

	birthDateTime, err := time.Parse("01/02/2006", birthDate)
	if err != nil {
		panic(err)
	}

	patient := blockchain.NewPatient(firstName, lastName, useName, birthDateTime)
	prescription := blockchain.NewPrescriptionEvent(patient)
	bc.AddEvent(prescription)

	log.Println("saving blockchain")
	if err := bc.Save(blockchainPath); err != nil {
		log.Fatal(err)
	}

}
