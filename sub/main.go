package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sapvs/otelechonats/common"
)

var nc *nats.Conn

func init() {
	nc, _ = nats.Connect(common.NATS_URL)
}

func main() {
	sub, er := nc.SubscribeSync(common.NATS_SUBJECT)
	if er != nil {
		log.Fatal(er)
	}
	defer sub.Unsubscribe()

	for {
		if msg, err := sub.NextMsg(10 * time.Second); err != nil {
			log.Fatal(err)
		} else {
			log.Printf("Message received %s", msg.Data)
		}

	}
	// if subs, err := nc.Subscribe(common.NATS_SUBJECT, func(msg *nats.Msg) {
	// 	log.Printf("receievd message %v", msg.Data)
	// }); err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	log.Println("SUbscribed")
	// 	defer subs.Unsubscribe()
	// }

}
