package blockchain

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestBlockchain(t *testing.T) {

	bc := NewBlockchain()

	birthDate, err := time.Parse("01/02/2006", "01/01/1970")
	if err != nil {
		t.Fatal(err)
	}

	patient := NewPatient("John", "Doe", "Doe", birthDate)

	bc.AddBlock(NewPrescriptionEvent(patient))

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Event: %v\n", block.Event)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}

}
