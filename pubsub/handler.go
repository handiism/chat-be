package pubsub

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/handiism/chat/user"
	"github.com/handiism/chat/utils"
)

type Handler struct {
	upgrader websocket.Upgrader
	service  user.Service
}

func NewHandler(service user.Service) Handler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return Handler{
		service:  service,
		upgrader: upgrader,
	}
}

var pubsub = &PubSub{}

func (h *Handler) WebSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if len(id) == 0 {
			c.JSON(http.StatusBadRequest, utils.NewFailResponse("User ID not provided", nil))
		}

		h.upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}

		conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.NewFailResponse(err.Error(), nil))
			return
		}
		defer conn.Close()

		user, err := h.service.FetchById(id)
		if err != nil {
			c.JSON(http.StatusNotFound, utils.NewFailResponse(err.Error(), nil))
			return
		}

		client := Client{
			ID:         user.Id.String(),
			Connection: conn,
		}

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				pubsub.RemoveClient(client)
				return
			}

			pubsub.ProcessMessage(client, messageType, p)
		}
	}
}
