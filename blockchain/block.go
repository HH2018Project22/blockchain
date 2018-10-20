package blockchain

import (
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	Timestamp     int64  `json:"ts"`
	Event         Event  `json:"evt"`
	PrevBlockHash []byte `json:"prv"`
	Hash          []byte `json:"hsh"`
	Nonce         int    `json:"nce"`
}

func (b *Block) Validate(bc *Blockchain) bool {
	pow := NewProofOfWork(b)
	return b.Event.Validate(bc) && pow.Validate()
}

func (b *Block) UnmarshalJSON(data []byte) error {

	var rawBlock map[string]*json.RawMessage
	if err := json.Unmarshal(data, &rawBlock); err != nil {
		return err
	}

	if err := json.Unmarshal(*rawBlock["ts"], &b.Timestamp); err != nil {
		return err
	}

	if err := json.Unmarshal(*rawBlock["prv"], &b.PrevBlockHash); err != nil {
		return err
	}

	if err := json.Unmarshal(*rawBlock["hsh"], &b.Hash); err != nil {
		return err
	}

	if err := json.Unmarshal(*rawBlock["nce"], &b.Nonce); err != nil {
		return err
	}

	var rawEvent map[string]*json.RawMessage
	if err := json.Unmarshal(*rawBlock["evt"], &rawEvent); err != nil {
		return err
	}

	var eventType EventType
	if err := json.Unmarshal(*rawEvent["typ"], &eventType); err != nil {
		return err
	}

	switch eventType {
	case Genesis:
		event := &GenesisEvent{}
		if err := json.Unmarshal(*rawBlock["evt"], event); err != nil {
			return err
		}
		b.Event = event
	case Prescription:
		event := &PrescriptionEvent{}
		if err := json.Unmarshal(*rawBlock["evt"], event); err != nil {
			return err
		}
		b.Event = event
	case Notification:
		event := &NotificationEvent{}
		if err := json.Unmarshal(*rawBlock["evt"], event); err != nil {
			return err
		}
		b.Event = event
	default:
		return fmt.Errorf("unknown event type '%s'", eventType)
	}

	return nil
}

func NewBlock(event Event, prevBlockHash []byte) *Block {

	block := &Block{
		Timestamp:     time.Now().Unix(),
		Event:         event,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock(NewGenesisEvent(), []byte{})
}
