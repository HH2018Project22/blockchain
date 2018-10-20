package blockchain

import "encoding/json"

type EventType string

type Event interface {
	json.Marshaler
	Type() EventType
	Data() []byte
	Validate(bc *Blockchain) bool
}
