package event

import (
	"fmt"
	"log"
)

type Dispatcher struct {
	jobs   chan job
	events map[Name]Listener
}

func NewDispatcher() *Dispatcher {
	d := &Dispatcher{
		jobs:   make(chan job),
		events: make(map[Name]Listener),
	}

	go d.consume()
	log.Println("Dispatcher initialized")
	return d
}

func (d *Dispatcher) Register(listener Listener, names ...Name) error {
	for _, name := range names {
		if _, ok := d.events[name]; ok {
			log.Printf("Event '%s' already registered", name)
			return fmt.Errorf("the '%s' event is already registered", name)
		}
		d.events[name] = listener
		log.Printf("Registered event: %s", name)
	}
	return nil
}

func (d *Dispatcher) Dispatch(name Name, event interface{}) error {
	if _, ok := d.events[name]; !ok {
		return fmt.Errorf("the '%s' event is not registered", name)
	}

	d.jobs <- job{eventName: name, eventType: event}
	log.Printf("Dispatched event: %s", name)

	return nil
}

func (d *Dispatcher) consume() {
	for job := range d.jobs {
		log.Printf("Consuming event: %s", job.eventName)
		d.events[job.eventName].Listen(job.eventType)
	}
}
