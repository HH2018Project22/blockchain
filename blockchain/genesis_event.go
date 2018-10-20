package blockchain

import "encoding/json"

const GenesisEventType EventType = "genesis"

type GenesisEvent struct{}

func (e *GenesisEvent) Type() EventType {
	return GenesisEventType
}

func (e *GenesisEvent) Validate(bc *Blockchain) error {
	return nil
}

func (e *GenesisEvent) Hash() []byte {
	return []byte(GenesisEventType)
}

func (e *GenesisEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type EventType `json:"type"`
	}{GenesisEventType})
}

func NewGenesisEvent() Event {
	return &GenesisEvent{}
}
