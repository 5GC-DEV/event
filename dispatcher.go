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
	go func() {
		fmt.Println("[DISPATCHER] Starting consume goroutine")
		d.consume()
	}()
	// go d.consume()
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
	fmt.Printf("*********** Dispatching Event Name: %s, Event Value: %#v\n", name, event)

	if _, ok := d.events[name]; !ok {
		return fmt.Errorf("the '%s' event is not registered", name)
	}

	// Print channel occupancy and capacity
	fmt.Printf("*************** Before assign Job Event Name: %s\n", name)
	fmt.Println("Channel occupancy:", len(d.jobs))
	fmt.Println("Channel capacity:", cap(d.jobs))

	// Send job to the channel
	d.jobs <- job{eventName: name, eventType: event}

	fmt.Printf("Dispatched event: %s\n", name)
	return nil
}

func (d *Dispatcher) consume() {
	fmt.Println("[DISPATCHER] consume loop started")
	for eventTask := range d.jobs {
		fmt.Printf("Consuming event: %s\n", eventTask.eventName)

		// Process each job in its own goroutine
		go func(j job) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Recovered in job handler: %v\n", r)
				}
			}()
			d.events[j.eventName].Listen(j.eventType)
		}(eventTask)
	}
	fmt.Println("consume End")
}
