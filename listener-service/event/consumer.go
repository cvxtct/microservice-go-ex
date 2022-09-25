package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

// to store received events from the queue
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// instance of the consumer
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	// random queue

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		// bind our channel to each of these topics
		ch.QueueBind(
			q.Name,
			s, // topic
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}
	// auto ignorage, is it exclusive, is it non local, no wait, arguments
	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	// listen until program termination, using a channel
	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)
			// TODO
			log.Println("To channel: ", q.Name, "With payload: ", payload)
			// fire a new go routine to make things faster
			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]", q.Name)

	// this channel makes running it "forever"
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		// authenticate

	// we can have as many cases as we want as long as we create the logic...

	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
