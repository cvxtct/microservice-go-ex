# Microservice Go

[Working with microservices Go](https://www.udemy.com/course/working-with-microservices-in-go/)

## Resources

### Packages

router package
go get github.com/go-chi/chi/v5
go get github.com/go-chi/chi/middleware
go get github.com/go-chi/cors          

mongo
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/mongo/options

mail
go get github.com/xhit/go-simple-mail/v2
go get github.com/vanng822/go-premailer/premailer

RabbitMQ
go get github.com/rabbitmq/amqp091-go

### Toolbox

[A simple example of how to create a reusable Go module with commonly used tools.](https://github.com/tsawler/toolbox)


## Describing cases

### HTTP + POSTGRES + MONGO + RABBITMQ

#### Case 1 - On frontend push "Test Broker" button

- Frontend creates the request and sends it to: http://localhost:8080 with the payload "empty post request" if no error it puts the broker's response to the screen along with the sent payload. As the request hitts the url, the routes route it and passes it to the Broker handler within handler.go. Broker function nothing but a responder with the message "Hit the broker" using writeJSON from helpers.go to send the response back to the caller. In the Broker function there is no logic implemented. 

- As an easy challenge, I have extended the Broker method to send the payload.Message to the MongoDB as well using the logRequest method to communicate with the logger-service.

#### Case 2 - On frontend push "Test Auth" button

- In this case the frontend will use the http://localhost:8080/handle endpoint and sends the payload populated with an email and password pair. Within broker-service the handler.go using handleSubmission method will process this request, and from this point every other functionalities entry point will be the handleSubmission method. Which is defined within the routes.go too. 
Since the request comming from the frontend has a field "action", the handleSubmission after reading the payload (readJSON method from the helpers.go) is able to make a decision regarding the function call. In this case since the "action" is "auth" the authenticate method will be called. 

- The authenticate method



 
