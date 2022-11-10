package main

import (
	"broker/internal/event"
	"broker/internal/logs"
	"broker/internal/types"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// If a new service with a new logic will rise
// Just add a new struct and extend the RequestPayload struct

// Broker does nothig extra but response the "Hit the broker" message
// if one calls the localhost:8080/
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := types.JsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	log.Println("Hit the broker")

	// log hit the broker http
	// err := app.logRequest("broker", payload.Message)
	// if err != nil {
	// 	app.errorJSON(w, err)
	// 	return
	// }

	_ = app.writeJSON(w, http.StatusOK, payload)

	// this not needed due to helpers.go
	// out, _ := json.MarshalIndent(payload, "", "\t")
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusAccepted)
	// w.Write(out)
}

// HandleSubmission actually the entrypoint of the services
// it is able to decide which handler to call based on the action field.
// clients must respect the API format:
//
//	{
//		action: "mail",
//		mail: {
//			from: "me@example.com",
//			to: "you@there.com",
//			subject: "Test Email",
//			message: "Hello there!",
//		}
//	}
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	// variable to store the request
	var requestPayload types.RequestPayload

	// extract the json from the requestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// switch to make the decision
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		// log via calling the logger service
		// app.logItem(w, requestPayload.Log)
		// log via calling the RabbitMQ listener service
		// app.logeventViaRabbit(w, requestPayload.Log)
		// Log item via RPC
		app.logItemViaRPC(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

// This will be not in use anymore, stay here for reference
// This was the http based logItem()
// logItem logs item received via http request
func (app *Config) logItem(w http.ResponseWriter, entry types.LogPayload) {
	// marshal entry (in prod MarshalIndent is not necessary)
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	// this is where we are going to send the request, logger-service endpoint
	logServiceURL := "http://logger-service/log"

	// define the request using the http module
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// needs a http header as well
	request.Header.Set("Content-Type", "application/json")

	// create a http client
	client := &http.Client{}

	// do the http request using the http client
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// clean after method execution
	defer response.Body.Close()

	// if logger-service response status code is not accepted
	// send back the error using errorJSON for the caller (in this case for the frontend)
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	// otherwise let's define a response
	var payload types.JsonResponse
	// fields are in this case:
	payload.Error = false
	payload.Message = "logged"
	// send back for the caller
	app.writeJSON(w, http.StatusAccepted, payload)
}

// This is not in use anymore, stay here for reference
// logeventViaRabbit puts the log payload into the rabbit queue
// using the specified channel
func (app *Config) logeventViaRabbit(w http.ResponseWriter, l types.LogPayload) {
	// simple as it is
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// again send back the response for the caller
	var payload types.JsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	app.writeJSON(w, http.StatusAccepted, payload)
}

// This is not in use anymore, stay here for reference
// pushToQueue will push the log message
func (app *Config) pushToQueue(name, msg string) error {
	// creates an event emitter from the event package
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	// create the log payload for the queue
	payload := types.LogPayload{
		Name: name,
		Data: msg,
	}

	// MarshalIndent no needed in production -> Marshal
	j, _ := json.MarshalIndent(&payload, "", "\t")
	// push it into the queue using "log-topic"
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}

	return nil

}

// Authenticate is juts an example logic how to pass information
// to a service.
func (app *Config) authenticate(w http.ResponseWriter, a types.AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()
	// make sure we gat back the correct statuscode

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService types.JsonResponse

	// decode the json from auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload types.JsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

// sendMail calls the mailer service and sends the Email payload
// steps are almost identical with the authenticate method
func (app *Config) sendMail(w http.ResponseWriter, msg types.MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// call the mailer service
	mailServiceUrl := "http://mailer-service/send"

	// post to mail service
	request, err := http.NewRequest("POST", mailServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer request.Body.Close()

	// make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mail service"))
		return
	}

	// send back json
	var payload types.JsonResponse
	payload.Error = false
	payload.Message = "Message sent to" + msg.To

	app.writeJSON(w, http.StatusAccepted, payload)
}

// define RPC payload struct
type RPCPayload struct {
	Name string
	Data string
}

// ligItemViaRPC will call the logger service and calls the remote function
func (app *Config) logItemViaRPC(w http.ResponseWriter, l types.LogPayload) {
	//
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// should convert l (type LogPayload) to RPCPayload instead of using struct literal
	// (S1016)go-staticcheck
	// rpcPayload := RPCPayload{
	// 	Name: l.Name,
	// 	Data: l.Data,
	// }

	rpcPayload := RPCPayload(l)

	var result string
	// the rpc method what we are going to call
	// it must be exported (start with capital)
	// remote function, payload, result
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Send the answer back to the frontend
	// This will store the info that the payload was processed
	payload := types.JsonResponse{
		Error:   false,
		Message: result,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}

// LogViaGRPC is using the compiled protobuf message and service definition
// to create remote call on the logger-service
func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	// payload variable to store the request payload
	var requestPayload types.RequestPayload

	// read the request payload into requestPayload variable
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// create the connection aka call the logger service gRPC server using no credentials
	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	// create a ew log service client with the connection object - this is from protobuf
	client := logs.NewLogServiceClient(conn)
	// create a context with timeout
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()

	// using contecxt call the WriteLog remote function and populate the protobuf message fields
	// the write log is at logger-service/cmd/api/grpc.go
	// confusing -> WriteLog sends back the "logged" message, however we define our own response message in the WriteLog!!??
	res, err := client.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload types.JsonResponse
	payload.Error = false
	// try out using the response
	// payload.Message = "logged"
	// this actually works
	payload.Message = res.String()

	app.writeJSON(w, http.StatusAccepted, payload)
}

// In order to be able to send a log request for the logger service
// Not part of the broker handler

// func (app *Config) logRequest(name, data string) error {
// 	var entry struct {
// 		Name string `json:"name"`
// 		Data string `json:"data"`
// 	}

// 	entry.Name = name
// 	entry.Data = data

// 	jsonData, _ := json.MarshalIndent(entry, "", "\t")
// 	logServiceURL := "http://logger-service/log"

// 	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return err
// 	}

// 	client := http.Client{}
// 	_, err = client.Do(request)
// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }
