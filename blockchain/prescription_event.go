package blockchain

import "encoding/json"

const Prescription EventType = "prescription"

type PrescriptionEvent struct{}

func (e *PrescriptionEvent) Type() EventType {
	return Prescription
}

func (e *PrescriptionEvent) Validate(bc *Blockchain) bool {
	return true
}

func (e *PrescriptionEvent) Data() []byte {
	return []byte{}
}

func (e *PrescriptionEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type EventType `json:"typ"`
	}{Prescription})
}

func NewPrescriptionEvent() Event {
	return &PrescriptionEvent{}
}
