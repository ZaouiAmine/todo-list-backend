package main

import (
	"encoding/json"
	"io"
	"time"

	"github.com/taubyte/go-sdk/event"
)

//export getTodos
func getTodos(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}
	setCORSHeaders(h)

	// Get room ID from query parameter
	roomID := getQueryParam(h, "room", "")
	if roomID == "" {
		h.Write([]byte("Missing room parameter"))
		h.Return(400)
		return 1
	}

	// Get all todos from database for this room
	todos, err := getAllTodos(roomID)
	if err != nil {
		return handleHTTPError(h, err, 500)
	}

	// Send response
	response := TodosResponse{
		Success: true,
		Message: "Todos retrieved successfully",
		Data:    todos,
	}

	return sendJSONResponse(h, response)
}

//export createTodo
func createTodo(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}
	setCORSHeaders(h)

	// Get room ID from query parameter
	roomID := getQueryParam(h, "room", "")
	if roomID == "" {
		h.Write([]byte("Missing room parameter"))
		h.Return(400)
		return 1
	}

	// Read request body
	body, err := io.ReadAll(h.Body())
	if err != nil {
		return handleHTTPError(h, err, 400)
	}

	// Parse request
	var req CreateTodoRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		h.Write([]byte("Invalid JSON"))
		h.Return(400)
		return 1
	}

	// Validate input
	if req.Text == "" {
		h.Write([]byte("Text is required"))
		h.Return(400)
		return 1
	}

	// Create new todo
	now := time.Now()
	todo := Todo{
		ID:        generateID(),
		RoomID:    roomID,
		Text:      req.Text,
		Completed: false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save to database
	err = saveTodo(todo)
	if err != nil {
		return handleHTTPError(h, err, 500)
	}

	// Publish real-time update
	publishTodoUpdate("created", &todo)

	// Send response
	response := TodoResponse{
		Success: true,
		Message: "Todo created successfully",
		Data:    &todo,
	}

	return sendJSONResponse(h, response)
}

//export updateTodo
func updateTodo(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}
	setCORSHeaders(h)

	// Get room ID and todo ID from query parameters
	roomID := getQueryParam(h, "room", "")
	if roomID == "" {
		h.Write([]byte("Missing room parameter"))
		h.Return(400)
		return 1
	}

	id := getQueryParam(h, "id", "")
	if id == "" {
		h.Write([]byte("Missing ID parameter"))
		h.Return(400)
		return 1
	}

	// Read request body
	body, err := io.ReadAll(h.Body())
	if err != nil {
		return handleHTTPError(h, err, 400)
	}

	// Parse request
	var req UpdateTodoRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		h.Write([]byte("Invalid JSON"))
		h.Return(400)
		return 1
	}

	// Get existing todo
	todo, err := getTodo(roomID, id)
	if err != nil {
		h.Write([]byte("Todo not found"))
		h.Return(404)
		return 1
	}

	// Update todo
	todo.Text = req.Text
	todo.Completed = req.Completed
	todo.UpdatedAt = time.Now()

	// Save to database
	err = saveTodo(*todo)
	if err != nil {
		return handleHTTPError(h, err, 500)
	}

	// Publish real-time update
	publishTodoUpdate("updated", todo)

	// Send response
	response := TodoResponse{
		Success: true,
		Message: "Todo updated successfully",
		Data:    todo,
	}

	return sendJSONResponse(h, response)
}

//export deleteTodoEndpoint
func deleteTodoEndpoint(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}
	setCORSHeaders(h)

	// Get room ID and todo ID from query parameters
	roomID := getQueryParam(h, "room", "")
	if roomID == "" {
		h.Write([]byte("Missing room parameter"))
		h.Return(400)
		return 1
	}

	id := getQueryParam(h, "id", "")
	if id == "" {
		h.Write([]byte("Missing ID parameter"))
		h.Return(400)
		return 1
	}

	// Get existing todo for real-time update
	todo, err := getTodo(roomID, id)
	if err != nil {
		h.Write([]byte("Todo not found"))
		h.Return(404)
		return 1
	}

	// Delete from database
	err = deleteTodo(roomID, id)
	if err != nil {
		return handleHTTPError(h, err, 500)
	}

	// Publish real-time update
	publishTodoUpdate("deleted", todo)

	// Send response
	response := TodoResponse{
		Success: true,
		Message: "Todo deleted successfully",
	}

	return sendJSONResponse(h, response)
}
