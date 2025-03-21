package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var redisClient *redis.Client

const (
	streamName  = "orders"
	groupName   = "order-consumer-group"
	consumerNum = 3
	producerNum = 2
)

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

// Producer simulates an order being placed
func startProducer(id int) {
	for {
		orderID := strconv.Itoa(rand.Intn(10000))
		customer := fmt.Sprintf("customer_%d", rand.Intn(100))

		args := &redis.XAddArgs{
			Stream: streamName,
			Values: map[string]interface{}{
				"order_id":  orderID,
				"customer":  customer,
				"timestamp": time.Now().Format(time.RFC3339),
			},
		}

		msgID, err := redisClient.XAdd(ctx, args).Result()
		if err != nil {
			log.Printf("Producer %d: failed to add order: %v", id, err)
		} else {
			log.Printf("Producer %d: added order %s by %s", id, msgID, customer)
		}

		time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
	}
}

// Consumer reads from Redis Stream
func startConsumer(consumerID int) {
	consumerName := fmt.Sprintf("consumer-%d", consumerID)

	for {
		streams, err := redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: consumerName,
			Streams:  []string{streamName, ">"},
			Count:    1,
			Block:    5 * time.Second,
		}).Result()

		if err != nil {
			if err == redis.Nil {
				continue
			}
			log.Printf("Consumer %d: error reading: %v", consumerID, err)
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				orderID := msg.Values["order_id"]
				customer := msg.Values["customer"]

				log.Printf("Consumer %d: processing order %v for %v", consumerID, orderID, customer)

				// Simulate processing time
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				id, _ := strconv.Atoi(orderID.(string))
				if id%2 == 0 {
					// Simulate message processing failure
					log.Printf("Consumer %d: failed to process message %s", consumerID, msg.ID)
				} else {
					// Acknowledge the message
					_, err := redisClient.XAck(ctx, streamName, groupName, msg.ID).Result()
					if err != nil {
						log.Printf("Consumer %d: error acknowledging message %s: %v", consumerID, msg.ID, err)
					} else {
						log.Printf("Consumer %d: acknowledged message %s", consumerID, msg.ID)
					}
				}
			}
		}
	}
}

func createConsumerGroup() {
	if err := redisClient.Do(ctx, "XINFO", "GROUPS", streamName).Err(); err == nil {
		return
	}

	if err := redisClient.XGroupCreateMkStream(ctx, streamName, groupName, "$").Err(); err != nil {
		log.Fatalf("Could not create consumer group: %v", err)
	}
}

func main() {
	initRedis()
	createConsumerGroup()

	log.Println("Starting Producers and Consumers...")

	for i := 0; i < producerNum; i++ {
		go startProducer(i + 1)
	}

	for i := 0; i < consumerNum; i++ {
		go startConsumer(i + 1)
	}

	select {} // Block forever
}
