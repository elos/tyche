package tyche

import (
	"container/heap"

	"github.com/elos/autonomous"
	"github.com/elos/models"
)

type Bid struct {
	Action    *models.Action
	Salience  float64
	Registree string
}

type Auction struct {
	autonomous.Life
	autonomous.Managed
	autonomous.Stopper

	Leaders chan *Bid
	*PriorityQueue
	Orders  chan *Bid
	Cancels chan string
	books   map[string]*Item
	leader  *Bid
}

func NewAuction() *Auction {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	return &Auction{
		Life:          autonomous.NewLife(),
		Stopper:       make(autonomous.Stopper),
		Leaders:       make(chan *Bid),
		PriorityQueue: &pq,
		Orders:        make(chan *Bid),
		Cancels:       make(chan string),
		books:         make(map[string]*Item),
	}
}

func (a *Auction) Start() {

Run:
	for {
		select {
		case b := <-a.Orders:
			apriori, ok := a.books[b.Registree]
			if !ok {
				apriori := &Item{
					value:    b,
					priority: b.Salience,
				}

				a.books[b.Registree] = apriori

				a.PriorityQueue.Push(apriori)

			} else {
				a.PriorityQueue.update(apriori, b, b.Salience)
			}
		case registree := <-a.Cancels:
			record, ok := a.books[registree]

			if !ok {
				return
			}

			heap.Remove(a.PriorityQueue, record.index)
			delete(a.books, registree)
		case <-a.Stopper:
			break Run
		}

		if max := a.PriorityQueue.max().value.(*Bid); max != a.leader {
			a.leader = max
			a.Leaders <- a.leader
		}
	}
}

func (a *Auction) Bid(as string, salience float64, action *models.Action) {
	a.Orders <- &Bid{
		Registree: as,
		Salience:  salience,
		Action:    action,
	}
}

func (a *Auction) Forfeit(as string) {
	a.Cancels <- as
}
