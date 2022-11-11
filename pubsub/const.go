package pubsub

import "github.com/gorilla/websocket"

const (
	publish     = "publish"
	subscribe   = "subscribe"
	unsubscribe = "unsubscribe"
)

type Subscription struct {
	Topic   string
	Clients *[]Client
}

type Client struct {
	ID         string
	Connection *websocket.Conn
}

type Message struct {
	Action  string `json:"action"`
	Topic   string `json:"topic"`
	Message string `json:"message"`
}
