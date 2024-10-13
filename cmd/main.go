package main

import (
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/raghav1030/go-ws-georedis/cmd/ws/location/handlers"
)

func main() {
	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if c.Get("host") == "localhost:3000" {
			c.Locals("Host", "Localhost:3000")
			return c.Next()
		}
		return c.Status(403).SendString("Request origin not allowed")
	})

	// Upgraded websocket request
	app.Get("/ws", websocket.New(handlers.RegisterWebsocketConnection))

	// ws://localhost:3000/ws
	log.Fatal(app.Listen(":3000"))
}
