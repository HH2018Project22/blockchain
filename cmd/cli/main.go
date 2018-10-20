package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/HH2018Project22/blockchain/blockchain"
)

var (
	blockchainPath = "bloodcoin.db"
)

func init() {
	flag.StringVar(&blockchainPath, "blockchain", blockchainPath, "Database file")
}

func main() {

	flag.Parse()

	if len(os.Args) < 2 {
		help()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {

	case "prescription":
		doPrescription(args)

	case "notification":
		doNotification(args)

	case "dump":
		doDump(args)

	default:
		help()
		os.Exit(1)
	}

}

func help() {
	fmt.Println("available commands: prescription, dump")
}

func getBlockchain() *blockchain.Blockchain {
	var bc *blockchain.Blockchain
	if _, err := os.Stat(blockchainPath); os.IsNotExist(err) {
		fmt.Println("creating new blockchain")
		bc = blockchain.NewBlockchain()
		if err = bc.Save(blockchainPath); err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("loading '%s'\n", blockchainPath)
		bc, err = blockchain.LoadBlockchain(blockchainPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	return bc
}
