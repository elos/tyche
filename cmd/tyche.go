package main

import (
	"log"
	"time"

	"github.com/elos/models"
	"github.com/elos/tyche"
)

func main() {
	auction := tyche.NewAuction()
	master := tyche.NewMaster(auction)

	events := make(tyche.Producer)
	master.AddProducer(events)

	sa := tyche.NewSleepAgent(100, master.Auction)
	master.AddConsumer(sa.Consumer())
	go master.StartAgent(sa)

	ra := tyche.NewReadingAgent(70, master.Auction)
	master.AddConsumer(ra.Consumer())
	go master.StartAgent(ra)

	go func() {
		for {
			select {
			case b := <-auction.Leaders:
				log.Print("NEW TOP BID: %+v", b)
			}
		}
	}()

	for i := 0; i < 200; i++ {
		time.Sleep(1 * time.Second)

		e := models.NewEvent()
		e.Name = "test "

		events <- e
	}
}
