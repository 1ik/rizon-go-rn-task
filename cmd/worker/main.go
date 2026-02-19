package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"rizon-test-task/internal/config"
	"rizon-test-task/internal/email"
	"rizon-test-task/internal/slack"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// .env file is optional - continue if it doesn't exist
		log.Println("Warning: .env file not found, using environment variables")
	}

	log.Println("Starting workers (email and Slack)...")

	// Load configurations
	emailCfg := config.GetEmailConfig()

	// Initialize senders
	emailSender := email.NewSMTPSender(emailCfg)
	slackSender := slack.NewMockSender() // Mock implementation that logs to console

	// Load RabbitMQ configuration
	rabbitMQCfg := config.GetRabbitMQConfig()

	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitMQCfg.URL())
	if err != nil {
		log.Fatal("failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Use WaitGroup to wait for both workers to finish
	var wg sync.WaitGroup
	wg.Add(2)

	// Start email worker
	go func() {
		defer wg.Done()
		if err := startEmailWorker(ctx, emailSender, conn); err != nil {
			log.Printf("Email worker error: %v", err)
		}
	}()

	// Start Slack worker
	go func() {
		defer wg.Done()
		if err := startSlackWorker(ctx, slackSender, conn); err != nil {
			log.Printf("Slack worker error: %v", err)
		}
	}()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for shutdown signal
	sig := <-sigChan
	log.Printf("Received signal: %v, shutting down workers...", sig)

	// Cancel context to signal workers to stop
	cancel()

	// Wait for both workers to finish
	wg.Wait()

	log.Println("All workers stopped")
}
