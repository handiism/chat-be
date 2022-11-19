package pubsub

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/handiism/chat/utils"
)

type Handler struct {
	upgrader websocket.Upgrader
}

func NewHandler() Handler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return Handler{
		upgrader: upgrader,
	}
}

var pubsub = &PubSub{}

func (h *Handler) WebSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}

		conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.NewFailResponse(err.Error(), nil))
			return
		}
		defer conn.Close()

		id, err := uuid.NewV4()

		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(err.Error(), nil))
		}
		client := Client{
			ID:         id.String(),
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
