package main

import (
	"encoding/base64"
	"flag"
	"log"

	"github.com/HH2018Project22/bloodcoin/blockchain"
)

var (
	notificationCommand = flag.NewFlagSet("notification", flag.ExitOnError)
	prescriptionHash    string
	notificationType    string
)

func init() {
	notificationCommand.StringVar(&prescriptionHash, "prescription", prescriptionHash, "Hash de la prescription")
	notificationCommand.StringVar(&notificationType, "type", notificationType, "Type de la notification")
}

func doNotification(args []string) {

	if err := notificationCommand.Parse(args); err != nil {
		panic(err)
	}

	bc := getBlockchain()

	log.Println("adding notification")

	prescriptionHashData, err := base64.StdEncoding.DecodeString(prescriptionHash)
	if err != nil {
		panic(err)
	}

	notification := blockchain.NewNotificationEvent(
		prescriptionHashData,
		blockchain.NotificationType(notificationType),
	)

	if err := bc.AddEvent(notification); err != nil {
		panic(err)
	}

	log.Println("saving blockchain")
	if err := bc.Save(blockchainPath); err != nil {
		panic(err)
	}

}
