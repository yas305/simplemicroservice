# simplemicroserviceReal-Time Microservices Application with Go, RabbitMQ, and WebSockets
Project Overview
This project demonstrates a real-time application using a microservices architecture in Go. It consists of two microservices that communicate asynchronously through RabbitMQ. The first microservice performs a basic calculation with input from a client, and the second microservice processes this result further. The final output is then pushed to the client in real-time using WebSockets.

Features
Microservice 1: Receives two numbers from the client, adds them, and sends the result to RabbitMQ.
Microservice 2: Consumes the result from RabbitMQ, doubles it, and sends it back to the client using WebSockets.
Real-time Updates: The client receives updates in real-time through an established WebSocket connection.
Technologies Used
Go: Efficient and straightforward language for building scalable microservices.
RabbitMQ: Robust message broker for handling asynchronous communication between services.
WebSockets: Enables real-time, bidirectional communication between the client and server.
HTML/JavaScript: For the client interface to interact with the microservices.
Getting Started
Prerequisites
Go (Version 1.x)
RabbitMQ Server
Git (for version control)
Installation
Clone the Repository

bash
Copy code
git clone https://github.com/yas305/simplemicroservice.git
cd your-repository
Start RabbitMQ Server

Ensure RabbitMQ is running on your machine. Refer to the RabbitMQ documentation for installation and running instructions.

Run the Microservices

Open two terminal windows, one for each microservice.

In the first terminal, navigate to the directory of the first microservice and run:

bash
Copy code
go run main.go
In the second terminal, navigate to the directory of the second microservice and run:

bash
Copy code
go run secondservice.go
Open the Client

Open index.html in your web browser to use the client interface.

Usage
Enter two numbers in the provided fields and click the 'Send' button. The result will be displayed in real-time as it is processed by the microservices.

