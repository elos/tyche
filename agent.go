package tyche

import (
	"github.com/elos/autonomous"
	"github.com/elos/data"
	"github.com/elos/models"
)

type Agent struct {
	autonomous.Life
	autonomous.Managed
	autonomous.Stopper

	identifier   data.ID
	bettingPower int
	auction      *Auction
	stream       Consumer
}

func NewAgent(name string, power int, auction *Auction) *Agent {
	return &Agent{
		Life:         autonomous.NewLife(),
		Stopper:      make(autonomous.Stopper),
		auction:      auction,
		stream:       make(Consumer),
		bettingPower: power,
		identifier:   data.ID(name),
	}
}

func (a *Agent) Consumer() Consumer {
	return a.stream
}

type SleepAgent struct {
	*Agent
}

func NewSleepAgent(power int, auction *Auction) *SleepAgent {
	return &SleepAgent{
		Agent: NewAgent("sleep", power, auction),
	}
}

func (a *SleepAgent) Start() {
	action := models.NewAction()
	action.Name = "Sleep Bid"
	bid := 0

Run:
	for {
		select {
		case <-a.stream:
			a.auction.Bid("sleep", float64(bid), action)
		case <-a.Stopper:
			break Run
		}

		if bid < a.bettingPower {
			bid += 5
		}
	}

	a.auction.Forfeit("sleep")
}

type ReadingAgent struct {
	*Agent
}

func NewReadingAgent(power int, auction *Auction) *ReadingAgent {
	return &ReadingAgent{
		Agent: NewAgent("reading", power, auction),
	}
}

func (a *ReadingAgent) Start() {
	action := models.NewAction()
	action.Name = "Read Bid"
	bid := 0

Run:
	for {
		select {
		case <-a.stream:
			a.auction.Bid("reading", float64(bid), action)
		case <-a.Stopper:
			break Run
		}

		if bid < a.bettingPower {
			bid += 10
		}
	}

	a.auction.Forfeit("reading")
}
