package ws

type Server struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

type Hub struct {
	Servers    map[string]*Server
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Servers:    make(map[string]*Server),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5), // Buffer size of 5
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			if server, ok := h.Servers[cl.ServerID]; ok {
				// Register client only if server exists
				if _, exists := server.Clients[cl.ID]; !exists {
					server.Clients[cl.ID] = cl
				}
			}
		case cl := <-h.Unregister:
			if server, ok := h.Servers[cl.ServerID]; ok {
				// Unregister client only if server exists
				if _, exists := server.Clients[cl.ID]; exists {
					// Notify others only if there are remaining clients
					if len(server.Clients) != 0 {
						h.Broadcast <- &Message{
							Content:  "User has left the chat",
							ServerID: cl.ServerID,
							Username: cl.Username,
						}
					}

					delete(server.Clients, cl.ID)
					close(cl.Message)
				}
			}

		case m := <-h.Broadcast:
			if server, ok := h.Servers[m.ServerID]; ok {
				// Broadcast message only if server exists
				for _, cl := range server.Clients {
					cl.Message <- m
				}
			}
		}
	}
}
