package tyche

import (
	"testing"
	"time"

	"github.com/elos/models"
)

func TestMerge(t *testing.T) {
	p1 := make(Producer)
	p2 := make(Producer)
	p3 := make(Producer)
	p4 := make(Producer)

	p12 := merge(p1, p2)
	p34 := merge(p3, p4)

	p := merge(p12, p34)

	go func() {
		var e *models.Event

		e = models.NewEvent()
		e.Name = "1"
		p1 <- e
		t.Log("Sent Event '1'")

		e = models.NewEvent()
		e.Name = "2"
		p2 <- e
		t.Log("Sent Event '2'")

		e = models.NewEvent()
		e.Name = "4"
		p4 <- e
		t.Log("Sent Event '4'")

		e = models.NewEvent()
		e.Name = "3"
		p3 <- e
		t.Log("Sent Event '3'")
	}()

	var e *models.Event

	select {
	case e = <-p:
		t.Log("Recieved Event '3'")
		if e.Name != "1" {
			t.Fatal("Expected Event '1' first")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timed out waiting for Event '1'")
	}

	e = <-p
	if e.Name != "2" {
		t.Log("Expected Event '2' second")
	}

	e = <-p
	if e.Name != "4" {
		t.Log("Expected Event '4' third")
	}

	e = <-p
	if e.Name != "3" {
		t.Log("Expected Event '3' fourth")
	}

}
