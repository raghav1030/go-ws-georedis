package location

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/raghav1030/go-ws-georedis/cmd/redis"
)

type Message struct {
	Action    string      `json:"action"` // "update_location" or "find_nearby"
	UserID    string      `json:"user_id"`
	Latitude  float64     `json:"latitude,omitempty"`
	Longitude float64     `json:"longitude,omitempty"`
	Radius    float64     `json:"radius,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type Client struct {
	ID                 string
	WebSocketConn      *websocket.Conn
	ReceiveMessageChan chan *Message
	Manager            *LocationManager
}

func NewClient(id string, ws *websocket.Conn, manager *LocationManager) *Client {
	return &Client{
		ID:                 id,
		WebSocketConn:      ws,
		ReceiveMessageChan: make(chan *Message),
		Manager:            manager,
	}
}

var Wg sync.WaitGroup

// WriteMessages handles receiving messages from the WebSocket
func (c *Client) WriteMessages() {
	defer func() {
		Wg.Done()
		c.Manager.UnSubscribeClientChan <- c
		_ = c.WebSocketConn.Close()
	}()

	for {
		_, msg, err := c.WebSocketConn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		receivedMessage := Message{}
		json.Unmarshal(msg, &receivedMessage)

		switch receivedMessage.Action {
		case "update_location":
			// Store location in Redis
			redisManager := redis_manager.GetRedisManager()
			err := redisManager.AddUserLocation(c.ID, receivedMessage.Latitude, receivedMessage.Longitude)
			if err != nil {
				fmt.Println("Error updating location:", err)
			}

		case "find_nearby":
			// Fetch nearby users
			redisManager := redis_manager.GetRedisManager()
			users, err := redisManager.GetNearbyUsers(c.ID, receivedMessage.Radius)
			if err != nil {
				fmt.Println("Error fetching nearby users:", err)
			}

			response := Message{
				Action: "nearby_users",
				Data:   users,
			}

			responseData, _ := json.Marshal(response)
			c.WebSocketConn.WriteMessage(websocket.TextMessage, responseData)
		}
	}
}
