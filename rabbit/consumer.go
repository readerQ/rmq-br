package rabbit

import (
	"fmt"
	"log"
)

type RmqConsumer struct {
	queue  string
	min    int
	max    int
	wait   bool
	conn   *RmqConnecton
	writer MessageWriter
}

func NewConsumer(queue string, min, max int, wait bool) *RmqConsumer {

	return &RmqConsumer{
		queue: queue,
		max:   max,
		min:   min,
		wait:  wait,
	}
}

func (rc *RmqConsumer) WithConnection(conn *RmqConnecton) *RmqConsumer {

	rc.conn = conn
	return rc
}

func (rc *RmqConsumer) WithWriter(wrt MessageWriter) *RmqConsumer {

	rc.writer = wrt
	return rc
}

// func (rc *RmqConsumer) chanCreate() error {

// 	// defer conn.Close() // todo
// 	err := rc.conn.Connect()
// 	if err!=nil{
// 		return err
// 	}

// 	ch, err := rc.conn.conn.Channel()
// 	if err!=nil{
// 		return err
// 	}
// 	defer ch.Close()

// 	q, err := ch. QueueDeclare(
// 	  "hello", // name
// 	  false,   // durable
// 	  false,   // delete when unused
// 	  false,   // exclusive
// 	  false,   // no-wait
// 	  nil,     // arguments
// 	)

// 	return err
// }

func (rc *RmqConsumer) Consume(min, max int) error {

	// defer conn.Close() // todo
	err := rc.conn.Connect()
	if err != nil {
		return fmt.Errorf("connection error: %s", err.Error())
	}

	ch, err := rc.conn.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel creation error: %s", err.Error())
	}

	msgs, err := ch.Consume(
		rc.queue, // queue
		"rmq-br", // consumer
		false,    // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)

	if err != nil {
		return fmt.Errorf("consume error: %s", err.Error())
	}

	forever := make(chan struct{})
	stop := make(chan struct{})

	count := 0
	go func() {

		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			msg := Message{
				Queue: rc.queue,
				Index: count,
				Body:  d.Body,
			}
			err := rc.writer.WriteMessage(msg)

			if err != nil {
				stop <- struct{}{}
			}

			ch.Ack(d.DeliveryTag, false)
			count++

			if rc.max > 0 && count >= rc.max {
				log.Printf("limit %d reached\n", rc.max)
				stop <- struct{}{}
				return
			}

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {
	case <-stop:
		{
			log.Println("stop consume by signal")
		}
	case <-forever:
		{

		}
	}

	return nil
}
