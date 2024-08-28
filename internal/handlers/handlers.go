package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsChan = make(chan WsPayload)
var clients = make(map[WsConnection]string)

var (
	views = jet.NewSet(
		jet.NewOSFileSystemLoader("./html"),
		jet.InDevelopmentMode(),
	)
	upgradeConnection = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type WsConnection struct {
	*websocket.Conn
}

type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

type WsPayload struct {
	Action   string       `json:"action"`
	Username string       `json:"username"`
	Message  string       `json:"message"`
	Conn     WsConnection `json:"-"`
}

func WsEndpoint(c *gin.Context) {
	ws, err := upgradeConnection.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	c.JSON(200, gin.H{"message": "Client connected to end point"})

	var response WsJsonResponse
	response.Message = `<em><small>Connected to server</small></em>`

	conn := WsConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}

	go ListenForWs(&conn)
}

func ListenForWs(conn *WsConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			log.Println(err)
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func ListenForWsChan() {
	var response WsJsonResponse

	for {
		e := <-wsChan
		response.Action = "Got here"
		response.Message = fmt.Sprintf("Some message and action %s", e.Action)
		broadcastToAll(response)
	}
}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println(err)
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func Home(ctx *gin.Context) {
	err := renderTemplate(ctx, "home.gohtml", nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func renderTemplate(c *gin.Context, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(c.Writer, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
