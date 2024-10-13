package location

type LocationManager struct {
	Clients                   []*Client
	SubscribeClientChan       chan *Client
	UnSubscribeClientChan     chan *Client
}

func (manager *LocationManager) Start() {
	for {
		select {
		case client := <-manager.SubscribeClientChan:
			manager.Clients = append(manager.Clients, client)

		case client := <-manager.UnSubscribeClientChan:
			for i, c := range manager.Clients {
				if c.ID == client.ID {
					manager.Clients = append(manager.Clients[:i], manager.Clients[i+1:]...)
				}
			}
		}
	}
}

var Manager = &LocationManager{
	Clients:               make([]*Client, 0),
	SubscribeClientChan:   make(chan *Client),
	UnSubscribeClientChan: make(chan *Client),
}
