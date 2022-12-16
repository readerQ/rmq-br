package rabbit

import (
	"github.com/rabbitmq/amqp091-go"
)

type RmqConnecton struct {
	Url string

	conn *amqp091.Connection
}

func (rconn *RmqConnecton) Connect() error {

	conn, err := amqp091.Dial(rconn.Url)
	rconn.conn = conn
	return err
}
