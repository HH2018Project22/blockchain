package main

import (
	"flag"
	"log"

	"github.com/davecgh/go-spew/spew"
)

var (
	listPrescriptionCommand = flag.NewFlagSet("list_prescription", flag.ExitOnError)
)

func init() {

}

func doListPrescription(args []string) {

	if err := listPrescriptionCommand.Parse(args); err != nil {
		panic(err)
	}

	bc := getBlockchain()

	log.Println("listing prescriptions")

	prescriptions := bc.ListPrescriptions()
	for _, p := range prescriptions {
		spew.Dump(p)
	}

}
