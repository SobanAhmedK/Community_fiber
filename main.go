package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

type WebSocketMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

var rooms = make(map[string][]*websocket.Conn)

func handleWebSocket(c *websocket.Conn) {
	room := c.Params("room")

	// Add connection to room
	rooms[room] = append(rooms[room], c)

	// Clean up on disconnect
	defer func() {
		for i, conn := range rooms[room] {
			if conn == c {
				rooms[room] = append(rooms[room][:i], rooms[room][i+1:]...)
				break
			}
		}
		c.Close()
	}()

	for {
		var msg WebSocketMessage
		if err := c.ReadJSON(&msg); err != nil {
			break
		}

		for _, conn := range rooms[room] {
			// /if conn != c {
			if err := conn.WriteJSON(msg); err != nil {
				break
			}
			// }
		}
	}
}

func main() {
	app := fiber.New()

	app.Static("/", "./public")

	app.Get("/ws/:room", websocket.New(handleWebSocket))
	app.Use(cors.New())

	app.Listen("192.168.56.1:3000")
}
