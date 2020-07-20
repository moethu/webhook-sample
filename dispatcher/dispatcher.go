package dispatcher

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
	"webhook-sample/issue"
)

// Dispatcher with single http client holding multiple hooks
type Dispatcher struct {
	Client *http.Client
	Hooks  map[string]string
	Mutex  *sync.Mutex
}

// NewDispatcher instanciates a new dispatcher object
func NewDispatcher() *Dispatcher {
	dispatcher := &Dispatcher{
		Client: &http.Client{},
		Hooks:  make(map[string]string),
		Mutex:  &sync.Mutex{},
	}
	return dispatcher
}

// Start starts the dispatcher loop sending data to hooks
func (d *Dispatcher) Start() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			d.Dispatch(issue.NewIssue())
		}
	}
}

// Add adds a webhook using name as a key
func (d *Dispatcher) Add(name string, hook string) {
	d.Mutex.Lock()
	log.Println("Adding:", name, "Webhook:", hook)
	d.Hooks[name] = hook
	d.Mutex.Unlock()
}

// Remove removes an existing webhook by key
func (d *Dispatcher) Remove(name string) {
	_, ok := d.Hooks[name]
	if ok {
		log.Println("Removing:", name)
		delete(d.Hooks, name)
	}
}

// Dispatch sends a serializable object to all registered endpoints
func (d *Dispatcher) Dispatch(data interface{}) {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()

	for name, hook := range d.Hooks {
		go func(name, hook string) {
			bytedata, err := json.Marshal(data)
			if err != nil {
				log.Println("Error mashalling json")
				return
			}

			request, err := http.NewRequest("POST", hook, bytes.NewReader(bytedata))
			if err != nil {
				log.Println("Error building request")
				return
			}

			response, err := d.Client.Do(request)
			if err != nil {
				log.Println("Error sending request")
				return
			}

			log.Println("Webhook:", hook, "Status:", response.StatusCode)
		}(name, hook)
	}
}
