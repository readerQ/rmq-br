package rabbit

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RmqSender struct {
	exchange             string
	routingKey           string
	mandatory, immediate bool

	conn   *RmqConnecton
	reader MessageReader
}

func NewSender(exchange string, routingKey string, mandatory, immediate bool) *RmqSender {

	return &RmqSender{
		exchange:   exchange,
		routingKey: routingKey,
	}
}

func (rc *RmqSender) WithConnection(conn *RmqConnecton) *RmqSender {

	rc.conn = conn
	return rc
}

func (rc *RmqSender) WithReader(rdr MessageReader) *RmqSender {

	rc.reader = rdr
	return rc
}

func (sndr *RmqSender) Send() (err error) {

	err = sndr.conn.Connect()
	defer sndr.conn.conn.Close()
	if err != nil {
		return fmt.Errorf("connection error: %s", err.Error())
	}

	ch, err := sndr.conn.conn.Channel()
	defer func() {
		_ = ch.Close()
	}()

	if err != nil {
		return fmt.Errorf("channel creation error: %s", err.Error())
	}
	count := 0
	defer func() { log.Printf("published %d messages", count) }()

	for {

		msg, found, err := sndr.reader.ReadMessage()

		if !found {

			log.Printf("no more files\n")
			return nil
		}

		if err != nil {
			log.Println(err)
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err = ch.PublishWithContext(ctx, sndr.exchange,
			sndr.routingKey,
			sndr.mandatory,
			sndr.immediate,
			amqp091.Publishing{
				ContentType: msg.ContentType,
				Body:        []byte(msg.Body),
			},
		)

		if err != nil {

			e2 := fmt.Errorf("publish error: %s", err.Error())
			log.Println(e2)
			return e2
		}

		count++
	}

}
