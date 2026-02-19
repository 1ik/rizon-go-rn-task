package message_broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"rizon-test-task/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	emailQueueName = "email_jobs"
	slackQueueName = "slack_jobs"
)

// EmailJob represents an email job message
type EmailJob struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// SlackJob represents a Slack job message
type SlackJob struct {
	UserEmail  string `json:"user_email"`
	DeviceID   string `json:"device_id"`
	Content    string `json:"content"`
	FeedbackID uint   `json:"feedback_id"`
}

// rabbitMQBroker implements MessageBroker using RabbitMQ
type rabbitMQBroker struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  *config.RabbitMQConfig
}

// NewRabbitMQBroker returns a RabbitMQ implementation of MessageBroker
func NewRabbitMQBroker(cfg *config.RabbitMQConfig) (MessageBroker, error) {
	url := cfg.URL()

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare email queue (idempotent - will create if doesn't exist)
	// Configure consumer timeout for delayed redelivery on failures
	_, err = channel.QueueDeclare(
		emailQueueName, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		amqp.Table{
			// Consumer timeout: 1 minute (60000 milliseconds)
			// If a message is not acked within this time, RabbitMQ will redeliver it
			"x-consumer-timeout": 60000, // 1 minute in milliseconds
		},
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare email queue: %w", err)
	}

	// Declare Slack queue (idempotent - will create if doesn't exist)
	_, err = channel.QueueDeclare(
		slackQueueName, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		amqp.Table{
			// Consumer timeout: 1 minute (60000 milliseconds)
			// If a message is not acked within this time, RabbitMQ will redeliver it
			"x-consumer-timeout": 60000, // 1 minute in milliseconds
		},
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare slack queue: %w", err)
	}

	log.Printf("RabbitMQ connected: %s (queues: %s, %s)", url, emailQueueName, slackQueueName)

	return &rabbitMQBroker{
		conn:    conn,
		channel: channel,
		config:  cfg,
	}, nil
}

// PublishEmailJob publishes an email job to the queue
func (r *rabbitMQBroker) PublishEmailJob(ctx context.Context, to, subject, body string) error {
	job := EmailJob{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	bodyBytes, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal email job: %w", err)
	}

	err = r.channel.PublishWithContext(
		ctx,
		"",              // exchange
		emailQueueName, // routing key (queue name)
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // Make message persistent
			ContentType:  "application/json",
			Body:         bodyBytes,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish email job: %w", err)
	}

	log.Printf("Published email job to queue: %s -> %s", emailQueueName, to)
	return nil
}

// Close closes the RabbitMQ connection
func (r *rabbitMQBroker) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			return fmt.Errorf("failed to close channel: %w", err)
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	}
	return nil
}

// PublishSlackJob publishes a Slack job to the queue
func (r *rabbitMQBroker) PublishSlackJob(ctx context.Context, userEmail, deviceID, content string, feedbackID uint) error {
	job := SlackJob{
		UserEmail:  userEmail,
		DeviceID:   deviceID,
		Content:    content,
		FeedbackID: feedbackID,
	}

	bodyBytes, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal slack job: %w", err)
	}

	err = r.channel.PublishWithContext(
		ctx,
		"",              // exchange
		slackQueueName,  // routing key (queue name)
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // Make message persistent
			ContentType:  "application/json",
			Body:         bodyBytes,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish slack job: %w", err)
	}

	log.Printf("Published slack job to queue: %s -> feedback ID: %d", slackQueueName, feedbackID)
	return nil
}

// GetEmailQueueName returns the email queue name (for workers)
func GetEmailQueueName() string {
	return emailQueueName
}

// GetSlackQueueName returns the Slack queue name (for workers)
func GetSlackQueueName() string {
	return slackQueueName
}
