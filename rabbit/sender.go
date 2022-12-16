package rabbit

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/readerQ/rmq-br/tools"
)

type RmqSender struct {
	exchange                    string
	routingKey                  string
	mandatory, immediate, quiet bool

	conn   *RmqConnecton
	reader MessageReader
}

func NewSender(exchange string, routingKey string, mandatory, immediate bool, quiet bool) *RmqSender {

	return &RmqSender{
		exchange:   exchange,
		routingKey: routingKey,
		mandatory:  mandatory,
		immediate:  immediate,
		quiet:      quiet,
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
	start := time.Now()
	var total uint64

	defer func() {
		speed := float64(total) / float64(time.Since(start).Seconds())
		log.Printf("published %d messages, %d bytes, avg %v ", count, total, tools.FormatNetSpeed(speed))
	}()

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

		if !sndr.quiet {
			var mot string
			if len(msg.Body) > 128 {
				mot = fmt.Sprintf("%s ... (%d bytes)", msg.Body[:128], len(msg.Body))
			} else {
				mot = string(msg.Body)
			}

			log.Printf("Send a message: %s", mot)
		}

		err = ch.PublishWithContext(ctx, sndr.exchange,
			sndr.routingKey,
			sndr.mandatory,
			sndr.immediate,
			amqp091.Publishing{
				ContentType: msg.ContentType,
				Body:        []byte(msg.Body),
			},
		)

		total += uint64(len(msg.Body))

		if err != nil {

			e2 := fmt.Errorf("publish error: %s", err.Error())
			log.Println(e2)
			return e2
		}

		count++
	}

}
