package handlers

import (
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/raghav1030/go-ws-georedis/cmd/ws/location"
)

func RegisterWebsocketConnection(c *websocket.Conn) {
    log.Println("WebSocket connection established")

    location.Wg.Add(1)
    fmt.Println("Creating client")

    client := location.NewClient(c.Params("userId"), c, location.Manager)
    // location.Manager.SubscribeClientChan <- client
    fmt.Println("Client created:", client)

    fmt.Println("Entering into WriteMessages goroutine")
    go func() {
        defer location.Wg.Done()  // Ensure the WaitGroup is decremented when done
        client.WriteMessages()
        fmt.Println("Exiting WriteMessages goroutine")
    }()

    fmt.Println("Waiting for goroutines to finish")
    location.Wg.Wait()  // This will block until the WaitGroup count is zero
    fmt.Println("All goroutines finished")
}
