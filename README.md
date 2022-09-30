# Microservice Go

[Working with microservices Go](https://www.udemy.com/course/working-with-microservices-in-go/)

    Remark: these notes are for practicing. Not all parts may being described.

- Easy to maintain / extend
- Robust
- Horisontal scaling
- Harder to write (each service is individual)

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

## Infrastructure

### Phase one

Each service run in docker container. Each service has a dockerfile to define image, the application directory and copy the binary to the work folder and a launch command. The build itself done by Makefile, then one docker-compuse.yml defines the services.


## Describing cases


### HTTP + POSTGRES + MONGO + RABBITMQ

#### **Case 1 - On frontend push "Test Broker" button**

Frontend creates the request and sends it to: http://localhost:8080 with the payload "empty post request" if no error it puts the broker's response to the screen along with the sent payload. As the request hitts the url, the routes route it and passes it to the Broker handler within handler.go. Broker function nothing but a responder with the message "Hit the broker" using writeJSON from helpers.go to send the response back to the caller. In the Broker function there is no logic implemented. 

As an easy challenge, I have extended the Broker method to send the payload.Message to the MongoDB as well using the logRequest method to communicate with the logger-service.

#### **Case 2 - On frontend push "Test Auth" button**

HTTP (without RabbitMQ)

    Note: this is not the way how we handle authentication in production

In this case the frontend will use the http://localhost:8080/handle endpoint and sends the payload populated with an email and password pair. Within broker-service the handler.go using handleSubmission method will get this request, and this will be the entry point. 

Since the request comming from the frontend has a field "action", the handleSubmission after reading the payload (readJSON method from the helpers.go (read json into the requestPayload strutct)) is able to make a decision regarding the function call. In this case since the "action" is "auth" the authenticate method will be called. 

The authenticate method has a receiver of Config and takes the http.ResponseWriter and a AuthPayload as argument. The AuthPayload is part of the RequestPayload which is sent by the frontend. 

```Go
type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
```

Stepps:
- marshall Payload
- call the service with http request -> request, err
- if err then send back the error msg and return
- create http client -> client
- do the request -> response, err and return
- if err then send the error back and return
- defer response body
- make sure we get back the correct status code
- Unauthroized -> send error and return
- Status not accepted -> send error and return
- create variable to read response body into
- decode json from auth service
- if err then send the error back and return
- if json from service body holds error then - send the error back and return 
- create variabla json response for response payload
- fill out its fields
- write json sends back the response for the frontend

Notice the places where error was handled: after sending the request to the service, after getting the response from the service, then while checking response body when decoding json, then during checking error within the response body. 

**Let's move to the authetication-service part:**

The pattern is the same, the main serves the service and the router routes the request to the appropriate handler. This service has a connection to the Postgresql database, for this reason it has a connectToDB and an openDB function to maintain the connection with the database. For CRUD methods the service has a /data/models.go [descriprion later].

During the request the Authenticatie method fires. It has a receiver of Config (this struct maintains the database connection). It implements the same patterns as others. Has a struct to store the request payload to work with, reads the response into this struct. 

To create a user object and retrieve the user from the database it calls form the model the corresponding method, in case of validate user app.Models.User.GetByEmail(requestPayload.Email) where app is the Config receiver which turns everything into one object, the Models is the main upper level struct which holds the User struct and the GetByEmail has the receiver pointer to User. 

Notice the way of object representation: app struct -> Models -> User -> method():

In main.go:
```Go
type Config struct {
	DB     *sql.DB
	Models data.Models // this is why Models is accessible from app
}
```
In models.go:
```Go
type Models struct {
	User User
}

// User is the structure which holds one user from the database.
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

Well, the Authenticate does the following stepps: 
- define a struct for the request payload
- reads the request into the requestPayload struct -> creates errorJSON if error and returns
- validates the user by email -> creats errorJSON if error or not valid and returns
- at this point it logs the authetication request using logRequest method and sending it to the log-service, thus new document being placed in the mongodb -> creates errorJSON if error and returns
- creates a payload with a jsonResponse json
- writes back the response for the broker (which then sends it towards the frontend)

#### **Case 3 - On frontend push "Test Log" button**

In this case the frontend creates the following:

```Go
 const payload = {
            action: "log",
            log: {
                name: "event",
                data: "Some kind of data",
            }
        }
```

Initially this request was sent as the aboves using http. Broker service receive it the broker's router routes it to the handleSubmission then it exctracts the action from the request which is in case "log" based on this the case will fire the logItem method which marshals the payload then creates the request for the log-service and send it to and writes back the log-service response for the frontend. Errors are handled at the same spots and way. 

At a certain point the RabbitMQ AMQP (Advanced Message Queue Protocol) was introduced to speed up communication and make it more conveninet. RabbitMQ runs in container, the listener-service implements the communication with this queue. 

The broker-service being extended with:
- events, consumer (not needed (yet at this point) despite the instructor replicated this file), emitter
- the main.go connects to the RabbitMQ, Config struct maintains the connection
- handler.go got logeventViaRabbit and pushToQueue methods 
- more simple, the fields are placed into the queue, a response back to the frontend that the action has been done
- pushToQueue creates a NeweventEmitter, then the emitter pushes the payload to the queue's (to the channel)


The emitter within the broker-service:
- Setup() to set up connection and returns declareExchange, this is the reason why the events.go is needed.
- NewEventEmitter() using the Setup() and returns the emitter

The listener-service listens and consumes each message from the channel. It creates a NewConsumer which has the listen method. 
The listen mechanism itself stands for the following steps:
- connects to the channel, if err return with
- defer channel
- declare random queue 
- we have one topic currently, however if there is more, a for loop pick and bind each channel to each topic
- initiate the message read from the channel
- creates a Go channel for the Go routine
- in Go routine a for loop starts in range of messages copy it into the Payload stuct, then another Go routine calls the handlePayload -> with this approach the listnere will be as fast as possible.

The handlePayload decides which event should call (using a switch case), in our case the topic is log_event -> logEvent() being called the log-service will be called. The rest is a usual service call. 

    We can have as many cases as we want as long as we create the logic to handle them...

### **Case 4 - On frontend push "Test Mail" button** 

# TODO

### Adding RPC to the logger-service

rpc.go added to the logger-service's main package with RPCServer and RPCPayload structs. 

LogInfo method is for putting the log message into mongo and send back an answer for the caler.

Still needs to run the RPC server and to achieve it within the main.go the followings are placed:
- rpcListen() method:
    - net.listen -> protocol, port
    - check for error
    - defer listen
    - infinite loop to listen and accept and serve connection within a go routine

- within main():
    - register the RPC server
    - listen in go routine

**-> Try out in rpc.go from line 37 - 53 : use data models method <-**

In broker-service, the handlers.go:
- got an RPCPayload struct which then will be populated with the LogPayload comming from the service that want's to log
- logItemViaRPC() method
- create client dialing the logger service on 5001 where the logger-service's rpc server is listen and serve
- call with the client -> RPCServer.LogInfo remote method


### Adding gRPC to the logger-service