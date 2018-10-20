package blockchain

import (
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	Timestamp     int64  `json:"timestamp"`
	Event         Event  `json:"event"`
	PrevBlockHash []byte `json:"previousBlockHash"`
	Hash          []byte `json:"hash"`
	Nonce         int    `json:"nonce"`
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

	if err := json.Unmarshal(*rawBlock["timestamp"], &b.Timestamp); err != nil {
		return err
	}

	if err := json.Unmarshal(*rawBlock["previousBlockHash"], &b.PrevBlockHash); err != nil {
		return err
	}

	if err := json.Unmarshal(*rawBlock["hash"], &b.Hash); err != nil {
		return err
	}

	if err := json.Unmarshal(*rawBlock["nonce"], &b.Nonce); err != nil {
		return err
	}

	var rawEvent map[string]*json.RawMessage
	if err := json.Unmarshal(*rawBlock["event"], &rawEvent); err != nil {
		return err
	}

	var eventType EventType
	if err := json.Unmarshal(*rawEvent["type"], &eventType); err != nil {
		return err
	}

	switch eventType {
	case GenesisEventType:
		event := &GenesisEvent{}
		if err := json.Unmarshal(*rawBlock["event"], event); err != nil {
			return err
		}
		b.Event = event
	case PrescriptionEventType:
		event := &PrescriptionEvent{}
		if err := json.Unmarshal(*rawBlock["event"], event); err != nil {
			return err
		}
		b.Event = event
	case NotificationEventType:
		event := &NotificationEvent{}
		if err := json.Unmarshal(*rawBlock["event"], event); err != nil {
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
