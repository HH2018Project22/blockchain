package blockchain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/btcsuite/btcutil/base58"
)

type Block struct {
	Timestamp     int64
	Event         Event
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func (b *Block) Validate(bc *Blockchain) error {
	pow := NewProofOfWork(b)
	if err := b.Event.Validate(bc); err != nil {
		return err
	}
	if err := pow.Validate(); err != nil {
		return err
	}
	return nil
}

func (b *Block) MarshalJSON() ([]byte, error) {
	rawBlock := map[string]interface{}{
		"timestamp":         b.Timestamp,
		"event":             b.Event,
		"previousBlockHash": base58.Encode(b.PrevBlockHash),
		"hash":              base58.Encode(b.Hash),
		"nonce":             b.Nonce,
	}
	return json.Marshal(rawBlock)
}

func (b *Block) UnmarshalJSON(data []byte) error {

	var rawBlock map[string]*json.RawMessage
	if err := json.Unmarshal(data, &rawBlock); err != nil {
		return err
	}

	if err := json.Unmarshal(*rawBlock["timestamp"], &b.Timestamp); err != nil {
		return err
	}

	var prevBlockHash string
	if err := json.Unmarshal(*rawBlock["previousBlockHash"], &prevBlockHash); err != nil {
		return err
	}
	b.PrevBlockHash = base58.Decode(prevBlockHash)

	var hash string
	if err := json.Unmarshal(*rawBlock["hash"], &hash); err != nil {
		return err
	}
	b.Hash = base58.Decode(hash)

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
