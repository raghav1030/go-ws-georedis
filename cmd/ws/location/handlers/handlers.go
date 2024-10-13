package handlers

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/raghav1030/go-ws-georedis/cmd/ws/location"
)

func RegisterWebsocketConnection( c *websocket.Conn) {
	location.Wg.Add(2)

	client := location.NewClient(c.Params("userId"), c, location.Manager)
	location.Manager.SubscribeClientChan <- client

	go client.WriteMessages()

	location.Wg.Wait()
}