package blockchain

import "encoding/json"

type EventType string

type Event interface {
	json.Marshaler
	Type() EventType
	Hash() []byte
	Validate(bc *Blockchain) error
}
