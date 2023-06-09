package repo

import (
	"context"
	"fmt"
	"nostr-ex/pkg/models"
	"time"

	eventUCase "nostr-ex/pkg/app/event/usecase"
	"nostr-ex/pkg/app/session/server/session"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

// func failOnError(err error, msg string) {
// 	if err != nil {
// 		log.Panicf("%s: %s", msg, err)
// 	}
// }

type Connector struct {
	url   string
	qName string

	metrics *Metrics
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
}

func NewConnector(url, qName string) *Connector {
	s := &Connector{
		url:     url,
		qName:   qName,
		metrics: NewMetrics("rabbit queue publisher"), //TODO:
	}

	return s
}

func (t *Connector) Connect() error {
	conn, err := amqp.Dial(t.url)
	if err != nil {
		return fmt.Errorf("%s %s", err, "Failed to connect to RabbitMQ")
	}
	t.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("%s %s", err, "Failed to open a channel")
	}
	t.channel = ch

	q, err := ch.QueueDeclare(
		t.qName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		return fmt.Errorf("%s %s", err, "Failed to declare a queue")
	}

	t.queue = &q
	return nil
}

func (t *Connector) DisConnect() {
	if t.conn == nil {
		return
	}

	if t.channel != nil {
		t.channel.Close()
	}

	t.conn.Close()

}

func (t *Connector) ConnectStatus() error {
	if t.channel == nil {
		return fmt.Errorf("channel == nil")
	}

	if t.queue == nil {
		return fmt.Errorf("queue == nil")
	}

	return nil
}

func (t *Connector) StartConsumer() error {

	err := t.ConnectStatus()
	if err != nil {
		return err
	}
	msgs, err := t.channel.Consume(
		t.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)

	if err != nil {
		return fmt.Errorf("channel.Consume error, %s", err.Error())
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Error("Connector StartConsumer() Error:", err)
			}
		}()
		ctx := context.Background()
		metrics := NewMetrics("rabbit queue consumer")
		eUCase := eventUCase.NewEventHandler()
		for d := range msgs {
			//TODO: delete log
			fmt.Printf("Received a message from MQ: %s\n", d.Body)
			shareMetrics := GetShareMetrics()
			shareMetrics.Dequeue(ctx)

			data := models.Event{
				SubID: "", //TODO:
				Data:  string(d.Body),
			}
			func() {
				t1 := time.Now()
				defer metrics.Duration(ctx, t1)
				err := eUCase.SaveEvent(&data)
				if err != nil {
					metrics.Fail(ctx)
					logrus.Error(err)
					return
				}
				//fmt.Printf("%#v\n", data)
				logrus.Debugf("%#v\n", data)
				session.ForEachSession(func(s session.SessionF) {
					s.OnDBDone()
				})

				metrics.Success(ctx)
				// mq := mqRepo.GetDBPublisher()
				// mq.Send(data.ID)
			}()
		}

		fmt.Println("Message queue consumer stop")
	}()

	return nil
}

func (t *Connector) Send(data []byte) error {

	t1 := time.Now()
	ctx := context.Background()
	defer t.metrics.Duration(ctx, t1)
	var err error
	defer func() {
		ctx := context.Background()
		if err != nil {
			t.metrics.Fail(ctx)
		} else {
			t.metrics.Success(ctx)
			metrics := GetShareMetrics()
			metrics.Enqueue(ctx)
		}
	}()

	err = t.ConnectStatus()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(),
		5*time.Second)
	defer cancel()

	err = t.channel.PublishWithContext(ctx,
		"",           // exchange
		t.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})
	if err != nil {
		return fmt.Errorf("%s %s", err, "Failed to publish a message")
	}

	fmt.Printf(" [E] Sent to MQ %s\n", data)

	return nil
}
