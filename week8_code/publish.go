package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

func main() {
	projectID := "msds434-mod8"
	topicID := "my-topic"
	msg := "Hello World"

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		fmt.Errorf("pubsub: NewClient: %w", err)
	}
	defer client.Close()

	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		fmt.Errorf("pubsub: result.Get: %w", err)
	}
	fmt.Printf("Published a message; msg ID: %v\n", id)

}
