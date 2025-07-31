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
		jobs:   make(chan job, 1),
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
	for job := range d.jobs {
		fmt.Printf("Consuming event: %s", job.eventName)
		d.events[job.eventName].Listen(job.eventType)
	}
}

/*func (d *Dispatcher) consume() {
	fmt.Println("[DISPATCHER] consume loop started")
	for job := range d.jobs {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[PANIC] while consuming event %s: %v\n", job.eventName, r)
				}
			}()
			fmt.Printf("Consuming event: %s, Value: %#v\n", job.eventName, job.eventType)
			listener, ok := d.events[job.eventName]
			if !ok {
				fmt.Printf("[ERROR] No listener for event: %s\n", job.eventName)
				return
			}
			listener.Listen(job.eventType)
		}()
	}
}
*/
