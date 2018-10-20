package main

import (
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

	notification := blockchain.NewNotificationEvent(
		[]byte(prescriptionHash),
		blockchain.NotificationType(notificationType),
	)

	if result := bc.AddEvent(notification); !result {
		panic("invalid notification")
	}

	log.Println("saving blockchain")
	if err := bc.Save(blockchainPath); err != nil {
		log.Fatal(err)
	}

}
