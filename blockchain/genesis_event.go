package blockchain

import "encoding/json"

const Genesis EventType = "genesis"

type GenesisEvent struct{}

func (e *GenesisEvent) Type() EventType {
	return Genesis
}

func (e *GenesisEvent) Validate(bc *Blockchain) bool {
	return true
}

func (e *GenesisEvent) Hash() []byte {
	return []byte(Genesis)
}

func (e *GenesisEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type EventType `json:"typ"`
	}{Genesis})
}

func NewGenesisEvent() Event {
	return &GenesisEvent{}
}
