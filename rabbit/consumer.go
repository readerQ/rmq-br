package rabbit

import (
	"fmt"
	"log"
	"time"
)

type RmqConsumer struct {
	queue string

	max    int
	wait   bool
	quiet  bool
	conn   *RmqConnecton
	writer MessageWriter
}

func NewConsumer(queue string, max int, wait bool, q bool) *RmqConsumer {

	return &RmqConsumer{
		queue: queue,
		max:   max,
		wait:  wait,
		quiet: q,
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

func (rc *RmqConsumer) Ping() (err error) {
	err = rc.conn.Connect()
	defer rc.conn.conn.Close()
	if err != nil {
		return fmt.Errorf("connection error: %s", err.Error())
	}

	ch, err := rc.conn.conn.Channel()
	defer func() {
		_ = ch.Close()
	}()

	if err != nil {
		return fmt.Errorf("channel creation error: %s", err.Error())
	}
	return nil
}

func (rc *RmqConsumer) Consume(min, max int) (err error) {

	err = rc.conn.Connect()
	defer rc.conn.conn.Close()
	if err != nil {
		return fmt.Errorf("connection error: %s", err.Error())
	}

	ch, err := rc.conn.conn.Channel()
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
		log.Printf("consumed %d messages, %d bytes, avg %.2f b/s ", count, total, speed)
	}()

	if rc.wait {
		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	}

	for {

		d, ok, err := ch.Get(rc.queue, false)

		if err != nil {
			log.Println(fmt.Errorf("error on get message: %s", err.Error()))
		}

		if !ok {
			if rc.wait {
				time.Sleep(250 * time.Millisecond)
				continue
			} else {
				if !rc.quiet {
					log.Println("no more messages in queue")
				}
				return nil
			}
		}

		if !rc.quiet {
			var mot string
			if len(d.Body) > 128 {
				mot = fmt.Sprintf("%s ... (%d bytes)", d.Body[:128], len(d.Body))
			} else {
				mot = string(d.Body)
			}

			log.Printf("Received a message: %s", mot)
		}

		total += uint64(len(d.Body))

		msg := Message{
			Queue: rc.queue,
			Index: count,
			Body:  d.Body,
		}

		err = rc.writer.WriteMessage(msg)

		if err != nil {
			return err
		}

		ch.Ack(d.DeliveryTag, false)
		count++

		if rc.max > 0 && count >= rc.max {
			log.Printf("max limit %d reached\n", rc.max)
			return nil
		}

	}

}
