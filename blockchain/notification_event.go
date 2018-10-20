package blockchain

import (
	"bytes"
	"encoding/json"
)

const Notification EventType = "notification"

const (
	Packaging  NotificationType = "packaging"
	Packaged   NotificationType = "packaged"
	Delivering NotificationType = "delivering"
	Delivered  NotificationType = "delivered"
	Transfused NotificationType = "transfused"
)

type NotificationType string

type NotificationEvent struct {
	PrescriptionHash []byte           `json:"pre"`
	NotificationType NotificationType `json:"not"`
}

func (e *NotificationEvent) Type() EventType {
	return Notification
}

func (e *NotificationEvent) Validate(bc *Blockchain) bool {
	return true
}

func (e *NotificationEvent) Hash() []byte {
	return bytes.Join([][]byte{
		[]byte(Notification),
		e.PrescriptionHash,
		[]byte(e.NotificationType),
	}, []byte{})
}

func (e *NotificationEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type             EventType        `json:"typ"`
		PrescriptionHash []byte           `json:"pre"`
		NotificationType NotificationType `json:"not"`
	}{Notification, e.PrescriptionHash, e.NotificationType})
}

func NewNotificationEvent(prescriptionHash []byte, notificationType NotificationType) Event {
	return &NotificationEvent{prescriptionHash, notificationType}
}
