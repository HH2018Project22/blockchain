package blockchain

import (
	"fmt"
	"strconv"
	"testing"
)

func TestBlockchain(t *testing.T) {

	bc := NewBlockchain()

	bc.AddBlock(NewPrescriptionEvent())
	bc.AddBlock(NewPrescriptionEvent())

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Event: %v\n", block.Event)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}

}
