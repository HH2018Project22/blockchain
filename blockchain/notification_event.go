package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"

	validator "gopkg.in/go-playground/validator.v9"
)

const NotificationEventType EventType = "notification"

const (
	Received   NotificationType = "received"
	Packaging  NotificationType = "packaging"
	Packaged   NotificationType = "packaged"
	Delivering NotificationType = "delivering"
	Delivered  NotificationType = "delivered"
	Transfused NotificationType = "transfused"
)

var notificationsSequence = []NotificationType{
	Received, Packaging, Packaged,
	Delivering, Delivered, Transfused,
}

func getParentNotificationType(nt NotificationType) NotificationType {
	var previous NotificationType
	for _, snt := range notificationsSequence {
		if snt == nt {
			return previous
		}
		previous = snt
	}
	return ""
}

type NotificationType string

type NotificationEvent struct {
	PrescriptionHash []byte           `json:"prescription" validate:"required"`
	NotificationType NotificationType `json:"notification" validate:"required"`
}

func (e *NotificationEvent) Type() EventType {
	return NotificationEventType
}

func (e *NotificationEvent) Validate(bc *Blockchain) bool {

	validate := validator.New()
	if err := validate.Struct(e); err != nil {
		fmt.Print(err)
		return false
	}

	prescription := bc.FindPrescription(e.PrescriptionHash)
	if prescription == nil {
		return false
	}

	events := bc.FindPrescriptionNotificationEvents(prescription.Hash())
	for _, ee := range events {
		if e.NotificationType == ee.NotificationType {
			return false
		}
	}

	if e.NotificationType == Received {
		return true
	}

	parentNotificationType := getParentNotificationType(e.NotificationType)

	for _, e := range events {
		if e.NotificationType == parentNotificationType {
			return true
		}
	}

	return false
}

func (e *NotificationEvent) Hash() []byte {
	return bytes.Join([][]byte{
		[]byte(NotificationEventType),
		e.PrescriptionHash,
		[]byte(e.NotificationType),
	}, []byte{})
}

func (e *NotificationEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type             EventType        `json:"type"`
		PrescriptionHash []byte           `json:"prescription"`
		NotificationType NotificationType `json:"notification"`
	}{NotificationEventType, e.PrescriptionHash, e.NotificationType})
}

func NewNotificationEvent(prescriptionHash []byte, notificationType NotificationType) Event {
	return &NotificationEvent{prescriptionHash, notificationType}
}
