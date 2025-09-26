package main

import (
	"encoding/json"
	"io"
	"time"

	"github.com/taubyte/go-sdk/event"
)

//export createRoom
func createRoom(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}
	setCORSHeaders(h)

	// Read request body
	body, err := io.ReadAll(h.Body())
	if err != nil {
		return handleHTTPError(h, err, 400)
	}

	// Parse request
	var req CreateRoomRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		h.Write([]byte("Invalid JSON"))
		h.Return(400)
		return 1
	}

	// Validate input
	if req.Name == "" {
		h.Write([]byte("Room name is required"))
		h.Return(400)
		return 1
	}

	// Create new room
	now := time.Now()
	room := Room{
		ID:        generateID(),
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save to database
	err = saveRoom(room)
	if err != nil {
		return handleHTTPError(h, err, 500)
	}

	// Send response
	response := RoomResponse{
		Success: true,
		Message: "Room created successfully",
		Data:    &room,
	}

	return sendJSONResponse(h, response)
}

//export getRoomEndpoint
func getRoomEndpoint(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}
	setCORSHeaders(h)

	// Get room ID from query parameter
	roomID := getQueryParam(h, "id", "")
	if roomID == "" {
		h.Write([]byte("Missing room ID parameter"))
		h.Return(400)
		return 1
	}

	// Get room from database
	room, err := getRoom(roomID)
	if err != nil {
		h.Write([]byte("Room not found"))
		h.Return(404)
		return 1
	}

	// Send response
	response := RoomResponse{
		Success: true,
		Message: "Room retrieved successfully",
		Data:    room,
	}

	return sendJSONResponse(h, response)
}
