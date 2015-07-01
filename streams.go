package tyche

import "github.com/elos/models"

// A Producer is a channel sending events
type Producer chan *models.Event

// A Consumer is a channel accepting events
type Consumer chan models.Event

// A Produceable is anything which can produce events
type Produceable interface {
	Producer() Producer
}

// A Consumerable is anything which can consume events
type Consumerable interface {
	Consumer() Consumer
}

// merge combines the events of two separate Producers
func merge(s1, s2 Producer) Producer {
	p := make(Producer)

	go func() {
		for {
			select {
			case e := <-s1:
				p <- e
			case e := <-s2:
				p <- e
			}
		}
	}()

	return p
}
