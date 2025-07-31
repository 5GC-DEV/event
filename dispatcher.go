package event

import (
	"fmt"
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
	fmt.Println("Dispatcher initialized")
	return d
}

func (d *Dispatcher) Register(listener Listener, names ...Name) error {
	for _, name := range names {
		if _, ok := d.events[name]; ok {
			fmt.Printf("Event '%s' already registered", name)
			return fmt.Errorf("the '%s' event is already registered", name)
		}
		d.events[name] = listener
		fmt.Printf("Registered event: %s", name)
	}
	return nil
}

func (d *Dispatcher) Dispatch(name Name, event interface{}) error {
	fmt.Printf("Dispatching Event Name: %s, Event Value: %#v\n", name, event)
	if _, ok := d.events[name]; !ok {
		return fmt.Errorf("the '%s' event is not registered", name)
	}

	d.jobs <- job{eventName: name, eventType: event}
	fmt.Printf("Dispatched event: %s", name)

	return nil
}

func (d *Dispatcher) consume() {
	fmt.Println("[DISPATCHER] consume loop started")
	for job := range d.jobs {
		fmt.Printf("Consuming event: %s", job.eventName)
		d.events[job.eventName].Listen(job.eventType)
	}
}
