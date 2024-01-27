package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Adjust the origin check as needed
    },
}

var clients = make(map[*websocket.Conn]bool) // connected clients

func handleConnections(w http.ResponseWriter, r *http.Request) {



    log.Println("Attempting to upgrade to WebSocket")
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade failed:", err)
        return
    }








    defer ws.Close()

    // Register new client
    clients[ws] = true

    for {
        // Read in a new message as JSON and map it to a Message object
        _, _, err := ws.ReadMessage()
        if err != nil {
            log.Printf("error: %v", err)
            delete(clients, ws)
            break
        }
        // Handle incoming messages (if needed)
    }
}


func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}

func main() {
	forever := make(chan bool)
    // Set up the HTTP server for WebSocket connections in a new goroutine
    go func() {
        http.HandleFunc("/ws", handleConnections)
        log.Println("WebSocket server starting on :9090")
        err := http.ListenAndServe(":9090", nil)
        if err != nil {
            log.Fatal("ListenAndServe: ", err)
        }
    }()


    

    // Set up RabbitMQ connection and channel
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    // Declare the RabbitMQ queue
    q, err := ch.QueueDeclare(
        "sumQueue", // queue name
        false,      // durable
        false,      // delete when unused
        false,      // exclusive
        false,      // no-wait
        nil,        // arguments
    )
    failOnError(err, "Failed to declare a queue")

    // Consume messages from the queue
    msgs, err := ch.Consume(
        q.Name, // queue
        "",     // consumer
        true,   // auto-ack
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
        nil,    // args
    )
    failOnError(err, "Failed to register a consumer")

    // Process messages from RabbitMQ
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
	
			// Convert the body to a string then parse as float64
			receivedNumber, err := strconv.ParseFloat(string(d.Body), 64)
			if err != nil {
				log.Printf("Error parsing number: %s", err)
				continue
			}
	
			// Double the number
			doubledNumber := receivedNumber * 2

	        wsResponse := map[string]float64{
				"originalResult": receivedNumber,
				"doubledResult":  doubledNumber,
			}
			// Send to every client connected via WebSocket
			for client := range clients {
				err := client.WriteJSON(wsResponse)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}()

    log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
  <-forever
}
