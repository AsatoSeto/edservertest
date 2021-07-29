package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
)

type Context struct {
	Title string
	Name  string
	Count int
}
type Status struct {
	System int `json:"Sys"`
	Weapon int `json:"Wep"`
	Engine int `json:"Eng"`
}

func brokerFunc(broker *Broker, stat chan []byte) {
	for {
		data := <-stat
		broker.Notifier <- data
	}
}

func main() {

	broker := NewServer()
	stat := make(chan []byte)

	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	e.Renderer = t
	e.GET("/hello", Hello)
	e.GET("/eventTest", func(c echo.Context) error {
		broker.ServeHTTP(c.Response(), c.Request())
		return nil
	})
	e.POST("/status", func(c echo.Context) error {
		b := make([]byte, 1024)
		i, _ := c.Request().Body.Read(b)
		s := Status{}
		json.Unmarshal(b[:i], &s)
		log.Println(s)
		body, _ := json.Marshal(s)
		stat <- body
		return nil
	})
	go brokerFunc(broker, stat)
	e.Start(":80")
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.Execute(w, Context{})
}

func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", nil)
}

// Example SSE server in Golang.
//     $ go run sse.go

type Broker struct {

	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan []byte

	// New client connections
	newClients chan chan []byte

	// Closed client connections
	closingClients chan chan []byte

	// Client connections registry
	clients map[chan []byte]bool
}

func NewServer() (broker *Broker) {
	// Instantiate a broker
	broker = &Broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()

	return
}

func (broker *Broker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// Make sure that the writer supports flushing.
	//
	flusher, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// Each connection registers its own message channel with the Broker's connections registry
	messageChan := make(chan []byte)

	// Signal the broker that we have a new connection
	broker.newClients <- messageChan

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		broker.closingClients <- messageChan
	}()

	// Listen to connection close and un-register messageChan
	// notify := rw.(http.CloseNotifier).CloseNotify()
	notify := req.Context().Done()

	go func() {
		<-notify
		broker.closingClients <- messageChan
	}()

	for {

		// Write to the ResponseWriter
		// Server Sent Events compatible
		fmt.Fprintf(rw, "data: %s\n\n", <-messageChan)

		// Flush the data immediatly instead of buffering it for later.
		flusher.Flush()
	}

}

func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.newClients:

			// A new client has connected.
			// Register their message channel
			broker.clients[s] = true
			log.Printf("Client added. %d registered clients", len(broker.clients))
		case s := <-broker.closingClients:

			// A client has dettached and we want to
			// stop sending them messages.
			delete(broker.clients, s)
			log.Printf("Removed client. %d registered clients", len(broker.clients))
		case event := <-broker.Notifier:

			// We got a new event from the outside!
			// Send event to all connected clients
			for clientMessageChan, _ := range broker.clients {
				clientMessageChan <- event
			}
		}
	}

}
