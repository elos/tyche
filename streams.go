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

// merge combines the events of separate Producers
func merge(streams ...Producer) Producer {
	p := make(Producer)

	for _, s := range streams {
		go func(stream, master Producer) {
			for e := range stream {
				master <- e
			}
		}(s, p)
	}

	return p
}
