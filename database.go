package main

import (
	"encoding/json"
	"fmt"

	"github.com/taubyte/go-sdk/database"
)

// getDatabaseConnection gets a database connection
func getDatabaseConnection() (database.Database, error) {
	return database.New("/todos")
}

// getRoomDatabaseConnection gets a database connection for rooms
func getRoomDatabaseConnection() (database.Database, error) {
	return database.New("/rooms")
}

// saveTodo saves a todo to the database
func saveTodo(todo Todo) error {
	db, err := getDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Convert todo to JSON
	todoJSON, err := json.Marshal(todo)
	if err != nil {
		return err
	}

	// Save to database with key pattern: /todos/{roomId}/{id}
	key := fmt.Sprintf("/todos/%s/%s", todo.RoomID, todo.ID)
	return db.Put(key, todoJSON)
}

// getTodo retrieves a todo from the database
func getTodo(roomID, id string) (*Todo, error) {
	db, err := getDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Get todo from database
	key := fmt.Sprintf("/todos/%s/%s", roomID, id)
	todoJSON, err := db.Get(key)
	if err != nil {
		return nil, err
	}

	// Parse JSON to Todo struct
	var todo Todo
	err = json.Unmarshal(todoJSON, &todo)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

// getAllTodos retrieves all todos from the database for a specific room
func getAllTodos(roomID string) ([]Todo, error) {
	db, err := getDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// List all todo keys for this room
	keys, err := db.List(fmt.Sprintf("/todos/%s/", roomID))
	if err != nil {
		return nil, err
	}

	var todos []Todo
	for _, key := range keys {
		// Get each todo
		todoJSON, err := db.Get(key)
		if err != nil {
			continue // Skip if can't get this todo
		}

		// Parse JSON to Todo struct
		var todo Todo
		err = json.Unmarshal(todoJSON, &todo)
		if err != nil {
			continue // Skip if can't parse this todo
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

// deleteTodo deletes a todo from the database
func deleteTodo(roomID, id string) error {
	db, err := getDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Delete todo from database
	key := fmt.Sprintf("/todos/%s/%s", roomID, id)
	return db.Delete(key)
}

// Room database operations

// saveRoom saves a room to the database
func saveRoom(room Room) error {
	db, err := getRoomDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Convert room to JSON
	roomJSON, err := json.Marshal(room)
	if err != nil {
		return err
	}

	// Save to database with key pattern: /rooms/{id}
	key := fmt.Sprintf("/rooms/%s", room.ID)
	return db.Put(key, roomJSON)
}

// getRoom retrieves a room from the database
func getRoom(id string) (*Room, error) {
	db, err := getRoomDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Get room from database
	key := fmt.Sprintf("/rooms/%s", id)
	roomJSON, err := db.Get(key)
	if err != nil {
		return nil, err
	}

	// Parse JSON to Room struct
	var room Room
	err = json.Unmarshal(roomJSON, &room)
	if err != nil {
		return nil, err
	}

	return &room, nil
}
