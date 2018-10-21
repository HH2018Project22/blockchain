package main

import (
	"flag"
	"fmt"
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

	var syncHook blockchain.BlockHookFunc
	if peerEndpoint != "" {
		syncHook = blockchain.CreateBlockSyncHook(peerEndpoint)
	}

	if _, err := os.Stat(blockchainPath); os.IsNotExist(err) {

		bc = blockchain.NewBlockchain(syncHook)
		if err = bc.Save(blockchainPath); err != nil {
			panic(err)
		}

	} else {

		bc, err = blockchain.LoadBlockchain(blockchainPath, syncHook)
		if err != nil {
			panic(err)
		}

	}

	return bc
}
