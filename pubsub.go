package main

import (
	"encoding/json"
	"fmt"
	"time"

	pubsub "github.com/taubyte/go-sdk/pubsub/node"
)

// publishTodoUpdate publishes a todo update to the real-time channel
func publishTodoUpdate(action string, todo *Todo) error {
	// Create room-specific channel for todo updates
	channelName := fmt.Sprintf("todo-updates-%s", todo.RoomID)
	channel, err := pubsub.Channel(channelName)
	if err != nil {
		return err
	}

	// Create update message
	update := map[string]interface{}{
		"action": action,
		"todo":   todo,
		"timestamp": fmt.Sprintf("%d", todo.UpdatedAt.Unix()),
	}

	// Convert to JSON
	message, err := json.Marshal(update)
	if err != nil {
		return err
	}

	// Publish the update
	return channel.Publish(message)
}

// publishTodoListUpdate publishes when the todo list changes for a specific room
func publishTodoListUpdate(roomID string) error {
	// Create room-specific channel for todo list updates
	channelName := fmt.Sprintf("todo-list-updates-%s", roomID)
	channel, err := pubsub.Channel(channelName)
	if err != nil {
		return err
	}

	// Create list update message
	update := map[string]interface{}{
		"action": "list-updated",
		"roomId": roomID,
		"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	}

	// Convert to JSON
	message, err := json.Marshal(update)
	if err != nil {
		return err
	}

	// Publish the update
	return channel.Publish(message)
}
