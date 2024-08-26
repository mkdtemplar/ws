package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

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

type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

func WsEndpoint(c *gin.Context) {
	ws, err := upgradeConnection.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	c.JSON(200, gin.H{"message": "Client connected to end point"})

	var response WsJsonResponse
	response.Message = `<em><small>Connected to server</small></em>`

	err = ws.WriteJSON(response)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
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
