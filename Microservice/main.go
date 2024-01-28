

package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "github.com/rs/cors"
    amqp "github.com/rabbitmq/amqp091-go"
)

// Define structure for incoming request
type RequestData struct {
    Number1 float64 `json:"number1"`
    Number2 float64 `json:"number2"`
}

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}


func main() {
    // Define the endpoint handler
   
    http.HandleFunc("/calculate", calculateHandler)


    // Set up CORS
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://127.0.0.1:5500"}, // Allow your client origin
        AllowedMethods: []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type"},
    })


    // Wrap the default mux with CORS middleware
 
    handler := c.Handler(http.DefaultServeMux)

    // Start the server
    fmt.Println("Server is running on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", handler))


}





func calculateHandler(w http.ResponseWriter, r *http.Request) {
    var data RequestData
    log.Printf("calculateHandler triggered")
    // Decode the request body
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        log.Printf("JSON decode error: %v", err)
        return
    }



    // Calculate the result
    result := data.Number1 + data.Number2
    fmt.Printf("Calculated Result: %f\n", result)


    // Respond to the client with the calculated result
    response := map[string]float64{"resultFromFirstMicroservice": result}
    json.NewEncoder(w).Encode(response)



    

    // Connect to RabbitMQ
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "sumQueue", // queue name
        false,      // durable
        false,      // delete when unused
        false,      // exclusive
        false,      // no-wait
        nil,        // arguments
    )
    failOnError(err, "Failed to declare a queue")

    // Convert result to byte slice
    body := fmt.Sprintf("%f", result)
    err = ch.PublishWithContext(
        context.Background(), // Use context.Background() for a basic context
        "",                   // exchange
        q.Name,               // routing key (queue name)
        false,                // mandatory
        false,                // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(body),
        })
    failOnError(err, "Failed to publish a message")
    log.Printf(" [x] Sent %s", body)
    

}


