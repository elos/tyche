package tyche

import (
	"github.com/elos/autonomous"
	"github.com/elos/models"
)

type Master struct {
	autonomous.Manager
	Auction       *Auction
	stream        Producer
	registrations chan *Consumer
	consumers     map[*Consumer]bool
}

func NewMaster(a *Auction) *Master {
	h := autonomous.NewHub()
	go h.Start()
	h.WaitStart()

	go h.StartAgent(a)

	m := &Master{
		Manager:       h,
		Auction:       a,
		stream:        make(chan *models.Event),
		registrations: make(chan *Consumer),
		consumers:     make(map[*Consumer]bool),
	}

	go m.run()

	return m
}

func (m *Master) run() {
	for {
		select {
		case e := <-m.stream:
			for c, _ := range m.consumers {
				*c <- *e
			}
		case c := <-m.registrations:
			m.consumers[c] = true
		}
	}
}

func (m *Master) AddProducer(p Producer) {
	go func() {
		for {
			select {
			case e := <-p:
				m.stream <- e
			}
		}
	}()
}

func (m *Master) AddConsumer(consumer Consumer) {
	m.registrations <- &consumer
}
