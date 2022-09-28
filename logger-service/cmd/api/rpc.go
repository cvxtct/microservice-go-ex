package main

import (
	"context"
	"log"
	"log-service/data"
	"time"
)

// RPC Server
type RPCServer struct {
}

// Received payload
// Must be the same that we do expect as a log item
type RPCPayload struct {
	Name string
	Data string
}

// resp *stirng is for sending message back to service which called this
// notice that in this part pointer as a parameter, on the
// caler part pass by reference
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	// Notice that this part is not using the data models
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}

	// Try out
	// Using the data models to put the log entry into mongo
	//
	// event := data.LogEntry{
	// 	Name: payload.Name,
	// 	Data: payload.Data,
	// }

	// app := Config{
	// 	Models: data.Models{},
	// }

	// err := app.Models.LogEntry.Insert(event)
	// if err != nil {
	// 	log.Println("error writin to mongo", err)
	// 	return err
	// }

	*resp = "Processed payload via RPC: " + payload.Name
	return nil
}
