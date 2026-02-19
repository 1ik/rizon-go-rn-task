package main

import (
	"context"
	"encoding/json"
	"log"

	"rizon-test-task/internal/email"
	"rizon-test-task/internal/message_broker"

	amqp "github.com/rabbitmq/amqp091-go"
)

// startEmailWorker starts the email worker that consumes from the email queue.
func startEmailWorker(ctx context.Context, emailSender email.EmailSender, conn *amqp.Connection) error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	// Declare queue (same as publisher)
	queueName := message_broker.GetEmailQueueName()
	_, err = channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{
			// Consumer timeout: 1 minute (60000 milliseconds)
			// If a message is not acked within this time, RabbitMQ will redeliver it
			"x-consumer-timeout": 60000, // 1 minute in milliseconds
		},
	)
	if err != nil {
		return err
	}

	// Set QoS to process one message at a time
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return err
	}

	// Consume messages
	msgs, err := channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack (we'll ack manually)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	log.Printf("Email worker started, consuming from queue: %s", queueName)

	// Process messages until context is cancelled
	for {
		select {
		case <-ctx.Done():
			log.Println("Email worker shutting down...")
			return nil
		case msg, ok := <-msgs:
			if !ok {
				// Channel closed (connection lost)
				log.Println("Email worker: message channel closed")
				return nil
			}
			processEmailJob(ctx, emailSender, msg)
		}
	}
}

// processEmailJob processes a single email job.
func processEmailJob(ctx context.Context, emailSender email.EmailSender, msg amqp.Delivery) {
	var job message_broker.EmailJob
	if err := json.Unmarshal(msg.Body, &job); err != nil {
		log.Printf("Error: failed to unmarshal email job: %v", err)
		msg.Nack(false, false) // Reject and don't requeue
		return
	}

	log.Printf("Processing email job: %s -> %s", job.To, job.Subject)

	// Send email with timeout
	// emailCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	// defer cancel()

	// if err := emailSender.SendEmail(emailCtx, job.To, job.Subject, job.Body); err != nil {
	// 	log.Printf("Error: failed to send email to %s: %v", job.To, err)
	// 	log.Printf("Message will be redelivered after consumer timeout (1 minute)")
	// 	return
	// }

	// Acknowledge message
	if err := msg.Ack(false); err != nil {
		log.Printf("Error: failed to acknowledge message: %v", err)
		return
	}
	log.Printf("Successfully sent email to %s", job.To)
}
