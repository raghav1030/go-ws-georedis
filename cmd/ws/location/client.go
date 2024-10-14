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

// WriteMessages handles receiving messages from the WebSocket using two separate Go routines
func (c *Client) WriteMessages() {
	defer func() {
		Wg.Done()
		c.Manager.UnSubscribeClientChan <- c
		_ = c.WebSocketConn.Close()
	}()
	
	fmt.Println("Starting message handler")

	// Only one goroutine to read from the WebSocket connection
	for {
		_, msg, err := c.WebSocketConn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading WebSocket message:", err)
			break
		}

		receivedMessage := Message{}
		err = json.Unmarshal(msg, &receivedMessage)
		if err != nil {
			fmt.Println("Error unmarshalling WebSocket message:", err)
			continue
		}

		fmt.Println("Received message:", receivedMessage)

		// Dispatch based on the action type
		switch receivedMessage.Action {
		case "update_location":
			fmt.Println("Updating location for", c.ID)
			// Store location in Redis
			redisManager := redis_manager.GetRedisManager()
			err := redisManager.AddUserLocation(c.ID, receivedMessage.Latitude, receivedMessage.Longitude)
			if err != nil {
				fmt.Println("Error updating location:", err)
			}

		case "find_nearby":
			fmt.Println("Finding nearby people for", c.ID)
			// Fetch nearby users
			redisManager := redis_manager.GetRedisManager()
			users, err := redisManager.GetNearbyUsers(c.ID, receivedMessage.Radius)
			if err != nil {
				fmt.Println("Error fetching nearby users:", err)
				continue
			}

			// Send the list of nearby users to the client
			response := Message{
				Action: "nearby_users",
				Data:   users,
			}

			responseData, _ := json.Marshal(response)
			err = c.WebSocketConn.WriteMessage(websocket.TextMessage, responseData)
			if err != nil {
				fmt.Println("Error sending nearby users response:", err)
			}
		default:
			fmt.Println("Unknown action:", receivedMessage.Action)
		}
	}

	fmt.Println("Exiting message handler")
}

