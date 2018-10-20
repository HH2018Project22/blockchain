package main

import (
	"flag"
	"log"
	"os"

	"github.com/HH2018Project22/blockchain"
	"github.com/davecgh/go-spew/spew"
)

var (
	blockchainPath = "blockchain.db"
)

func init() {
	flag.StringVar(&blockchainPath, "blockchain", blockchainPath, "Database file")
}

func main() {

	flag.Parse()

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

	spew.Dump(bc)

}
