package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/HH2018Project22/bloodcoin/blockchain"
)

var (
	blockchainPath = "bloodcoin.db"
	peerEndpoint   = ""
)

func init() {
	flag.StringVar(&blockchainPath, "blockchain", blockchainPath, "Database file")
	flag.StringVar(&peerEndpoint, "peer", peerEndpoint, "Peer endpoint")
}

func main() {

	flag.Parse()

	if len(flag.Args()) < 1 {
		help()
		os.Exit(1)
	}

	command := flag.Args()[0]
	args := flag.Args()[1:]

	switch command {

	case "prescription":
		doPrescription(args)

	case "list_prescription":
		doListPrescription(args)

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
	fmt.Println("available commands: prescription, dump, list_prescription")
}

func getBlockchain() *blockchain.Blockchain {
	var bc *blockchain.Blockchain

	if _, err := os.Stat(blockchainPath); os.IsNotExist(err) {

		fmt.Println("creating new blockchain")
		bc = blockchain.NewBlockchain(beforeBlockAdd)
		if err = bc.Save(blockchainPath); err != nil {
			panic(err)
		}

	} else {

		fmt.Printf("loading '%s'\n", blockchainPath)
		bc, err = blockchain.LoadBlockchain(blockchainPath, beforeBlockAdd)
		if err != nil {
			panic(err)
		}

	}

	return bc
}

func beforeBlockAdd(block *blockchain.Block) error {

	if peerEndpoint == "" {
		return nil
	}

	url := fmt.Sprintf("%s/blocks/new", peerEndpoint)

	data, err := json.Marshal(block)
	if err != nil {
		return err
	}

	res, err := http.Post(url, "json/application", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return errors.New("could not propagate block")
	}

	return nil

}
