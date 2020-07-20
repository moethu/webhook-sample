package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"webhook-sample/dispatcher"
)

// Webhook payload
type Webhook struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func main() {
	// instanciate a new dispatcher holding webhooks
	dispatcher := dispatcher.NewDispatcher()

	// setup HTTP server
	srv := &http.Server{Addr: ":8000", Handler: http.DefaultServeMux}

	// webhook registration
	http.HandleFunc("/register", func(resp http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		var data Webhook
		err := decoder.Decode(&data)
		if err != nil {
			resp.WriteHeader(http.StatusBadRequest)
			return
		}
		if data.Url != "" {
			dispatcher.Add(data.Name, data.Url)
		}
		resp.WriteHeader(200)
	})

	// webhook removal
	http.HandleFunc("/unregister", func(resp http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		var data Webhook
		err := decoder.Decode(&data)
		if err != nil {
			resp.WriteHeader(http.StatusBadRequest)
			return
		}
		dispatcher.Remove(data.Name)
		resp.WriteHeader(200)
	})

	// webhook loopback test printing body
	http.HandleFunc("/test", func(resp http.ResponseWriter, req *http.Request) {
		bodyBytes, _ := ioutil.ReadAll(req.Body)
		log.Println("Incoming:", string(bodyBytes))
		resp.WriteHeader(200)
	})

	// start dispatching webhooks
	go dispatcher.Start()

	log.Println("Starting Webserver")
	err := srv.ListenAndServe()

	if err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
