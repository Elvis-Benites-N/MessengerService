package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

// CreateServerReq represents the JSON request for creating a new server.
type CreateServerReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateServer creates a new server and adds it to the hub.
func (h *Handler) CreateServer(c *gin.Context) {
	var req CreateServerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Check if the server ID already exists in the hub.
	if _, ok := h.hub.Servers[req.ID]; ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Server ID already exists"})
		return
	}

	// Create the new server and add it to the hub.
	h.hub.Servers[req.ID] = &Server{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	c.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// JoinServer handles WebSocket connection upgrades and client registration to a server.
func (h *Handler) JoinServer(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade WebSocket connection"})
		return
	}

	serverID := c.Param("serverId")
	clientID := c.Query("userId")
	username := c.Query("username")

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		ServerID: serverID,
		Username: username,
	}

	m := &Message{
		Content:  "A new user has joined the chat",
		ServerID: serverID,
		Username: username,
	}

	// Register the client and broadcast the message to all clients in the server.
	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(h.hub)
}

// ServerRes represents the JSON response for server information.
type ServerRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetServers returns a list of servers from the hub.
func (h *Handler) GetServers(c *gin.Context) {
	Servers := make([]ServerRes, 0)

	for _, r := range h.hub.Servers {
		Servers = append(Servers, ServerRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	c.JSON(http.StatusOK, Servers)
}

// ClientRes represents the JSON response for client information.
type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// GetClients returns a list of clients in a specific server from the hub.
func (h *Handler) GetClients(c *gin.Context) {
	var clients []ClientRes
	serverID := c.Param("serverId")

	// Check if the server ID exists in the hub.
	if _, ok := h.hub.Servers[serverID]; !ok {
		clients = make([]ClientRes, 0)
		c.JSON(http.StatusOK, clients)
	}

	// Iterate through the clients in the server and append them to the response.
	for _, client := range h.hub.Servers[serverID].Clients {
		clients = append(clients, ClientRes{
			ID:       client.ID,
			Username: client.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}
