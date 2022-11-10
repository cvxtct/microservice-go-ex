package config

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type AppConfig struct {
	Rabbit *amqp.Connection
}
