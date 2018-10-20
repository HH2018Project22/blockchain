package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
)

var dumpCommand = flag.NewFlagSet("dump", flag.ExitOnError)

func doDump() {

	bc := getBlockchain()

	fmt.Printf("------------\n")
	for _, b := range bc.Blocks() {
		fmt.Printf("Hash: %x\n", b.Hash)
		fmt.Printf("Prev. hash: %x\n", b.PrevBlockHash)
		data, err := json.Marshal(b.Event)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Event: %s\n", data)
		fmt.Printf("Valid: %s\n", strconv.FormatBool(b.Validate(bc)))
		fmt.Printf("------------\n")
	}

}
