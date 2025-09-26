package main

import "time"

// Todo represents a todo item
type Todo struct {
	ID          string    `json:"id"`
	RoomID      string    `json:"roomId"`
	Text        string    `json:"text"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreateTodoRequest represents the request to create a new todo
type CreateTodoRequest struct {
	Text string `json:"text"`
}

// CreateRoomRequest represents the request to create a new room
type CreateRoomRequest struct {
	Name string `json:"name"`
}

// Room represents a todo room
type Room struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// RoomResponse represents the response for room operations
type RoomResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    *Room  `json:"data,omitempty"`
}

// UpdateTodoRequest represents the request to update a todo
type UpdateTodoRequest struct {
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

// TodoResponse represents the response for todo operations
type TodoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    *Todo  `json:"data,omitempty"`
}

// TodosResponse represents the response for getting all todos
type TodosResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []Todo `json:"data"`
}
