package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"rizon-test-task/internal/message_broker"
	"rizon-test-task/internal/slack"

	amqp "github.com/rabbitmq/amqp091-go"
)

// startSlackWorker starts the Slack worker that consumes from the Slack queue.
func startSlackWorker(ctx context.Context, slackSender slack.SlackSender, conn *amqp.Connection) error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	// Declare queue (same as publisher)
	queueName := message_broker.GetSlackQueueName()
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

	log.Printf("Slack worker started, consuming from queue: %s", queueName)

	// Process messages until context is cancelled
	for {
		select {
		case <-ctx.Done():
			log.Println("Slack worker shutting down...")
			return nil
		case msg, ok := <-msgs:
			if !ok {
				// Channel closed (connection lost)
				log.Println("Slack worker: message channel closed")
				return nil
			}
			processSlackJob(ctx, slackSender, msg)
		}
	}
}

// processSlackJob processes a single Slack job.
func processSlackJob(ctx context.Context, slackSender slack.SlackSender, msg amqp.Delivery) {
	var job message_broker.SlackJob
	if err := json.Unmarshal(msg.Body, &job); err != nil {
		log.Printf("Error: failed to unmarshal slack job: %v", err)
		msg.Nack(false, false) // Reject and don't requeue
		return
	}

	log.Printf("Processing slack job: feedback ID %d from %s", job.FeedbackID, job.UserEmail)

	// Format Slack message
	message := formatSlackMessage(job)

	// Send to Slack with timeout
	slackCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := slackSender.SendMessage(slackCtx, message); err != nil {
		log.Printf("Error: failed to send slack message for feedback ID %d: %v", job.FeedbackID, err)
		log.Printf("Message will be redelivered after consumer timeout (1 minute)")
		return
	}

	// Acknowledge message
	if err := msg.Ack(false); err != nil {
		log.Printf("Error: failed to acknowledge message: %v", err)
		return
	}
	log.Printf("Successfully sent slack message for feedback ID %d", job.FeedbackID)
}

// formatSlackMessage formats the Slack message with feedback details.
func formatSlackMessage(job message_broker.SlackJob) string {
	return fmt.Sprintf(`New Feedback Submitted

*User Email:* %s
*Device ID:* %s
*Feedback ID:* %d
*Content:*
%s
*Submitted:* %s`,
		job.UserEmail,
		job.DeviceID,
		job.FeedbackID,
		job.Content,
		time.Now().Format(time.RFC3339),
	)
}
