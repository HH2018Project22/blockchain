package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/HH2018Project22/blockchain"
)

var (
	blockchainPath = "blockchain.db"
	command        = "dump"
)

func init() {
	flag.StringVar(&blockchainPath, "blockchain", blockchainPath, "Database file")
	flag.StringVar(&command, "cmd", command, "Command")
}

func main() {

	flag.Parse()

	bc := getBlockchain()

	switch command {
	case "dump":
		for _, b := range bc.Blocks() {
			fmt.Printf("Prev. hash: %x\n", b.PrevBlockHash)
			fmt.Printf("Event: %v\n", b.Event)
			fmt.Printf("Hash: %x\n", b.Hash)
			fmt.Printf("Valid: %s\n", strconv.FormatBool(b.Validate(bc)))
		}
	case "prescription":
		log.Println("Adding prescription")
		prescription := blockchain.NewPrescriptionEvent(blockchain.Patient{
			FirstName: "John",
			LastName:  "Doe",
			UseName:   "Doe",
		})
		bc.AddBlock(prescription)
		if err := bc.Save(blockchainPath); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown command: '%s'", command)
	}

}

func getBlockchain() *blockchain.Blockchain {
	var bc *blockchain.Blockchain
	if _, err := os.Stat(blockchainPath); os.IsNotExist(err) {
		log.Println("creating new blockchain")
		bc = blockchain.NewBlockchain()
		if err = bc.Save(blockchainPath); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("loading blockchain '%s'", blockchainPath)
		bc, err = blockchain.LoadBlockchain(blockchainPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	return bc
}
