package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rizon-test-task/internal/config"
	"rizon-test-task/internal/email"
	"rizon-test-task/internal/message_broker"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("Starting email worker...")

	// Load email configuration
	emailCfg := config.GetEmailConfig()

	// Initialize email sender
	emailSender := email.NewSMTPSender(emailCfg)

	// Load RabbitMQ configuration
	rabbitMQCfg := config.GetRabbitMQConfig()

	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitMQCfg.URL())
	if err != nil {
		log.Fatal("failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal("failed to open channel:", err)
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
		nil,       // arguments
	)
	if err != nil {
		log.Fatal("failed to declare queue:", err)
	}

	// Set QoS to process one message at a time
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatal("failed to set QoS:", err)
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
		log.Fatal("failed to register consumer:", err)
	}

	log.Printf("Email worker started, consuming from queue: %s", queueName)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Process messages
	go func() {
		for msg := range msgs {
			processEmailJob(context.Background(), emailSender, msg)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Println("Shutting down email worker...")
}

func processEmailJob(ctx context.Context, emailSender email.EmailSender, msg amqp.Delivery) {
	var job message_broker.EmailJob
	if err := json.Unmarshal(msg.Body, &job); err != nil {
		log.Printf("Error: failed to unmarshal email job: %v", err)
		msg.Nack(false, false) // Reject and don't requeue
		return
	}

	log.Printf("Processing email job: %s -> %s", job.To, job.Subject)

	// Send email with timeout
	emailCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := emailSender.SendEmail(emailCtx, job.To, job.Subject, job.Body); err != nil {
		log.Printf("Error: failed to send email to %s: %v", job.To, err)
		msg.Nack(false, true) // Reject and requeue for retry
		return
	}

	// Acknowledge message
	if err := msg.Ack(false); err != nil {
		log.Printf("Error: failed to acknowledge message: %v", err)
	} else {
		log.Printf("Successfully sent email to %s", job.To)
	}
}
