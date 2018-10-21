package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/btcsuite/btcutil/base58"
)

var dumpCommand = flag.NewFlagSet("dump", flag.ExitOnError)

func doDump(args []string) {

	if err := dumpCommand.Parse(args); err != nil {
		log.Fatal(err)
	}

	bc := getBlockchain()

	fmt.Printf("------------\n")
	for _, b := range bc.Blocks() {
		fmt.Printf("Hash: %s\n", base58.Encode(b.Hash))
		fmt.Printf("Prev. hash: %s\n", base58.Encode(b.PrevBlockHash))
		data, err := json.Marshal(b.Event)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Event: %s\n", data)
		fmt.Printf("------------\n")
	}

}
