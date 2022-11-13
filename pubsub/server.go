package pubsub

import (
	"encoding/json"
)

type PubSub struct {
	Subscriptions []Subscription
}

func (s *PubSub) Send(client *Client, message string) {
	client.Connection.WriteMessage(1, []byte(message))
}

func (s *PubSub) RemoveClient(client Client) {
	for _, sub := range s.Subscriptions {
		for i := 0; i < len(*sub.Clients); i++ {
			if client.ID == (*sub.Clients)[i].ID {
				if i == len(*sub.Clients)-1 {
					*sub.Clients = (*sub.Clients)[:len(*sub.Clients)-1]
				} else {
					*sub.Clients = append((*sub.Clients)[:i], (*sub.Clients)[i+1:]...)
					i--
				}
			}
		}
	}
}

func (s *PubSub) ProcessMessage(client Client, messageType int, payload []byte) *PubSub {
	m := Message{}
	if err := json.Unmarshal(payload, &m); err != nil {
		s.Send(&client, "Server: Invalid payload")
	}

	switch m.Action {
	case publish:
		s.Publish(client, m.Topic, []byte(m.Message))
	case subscribe:
		s.Subscribe(&client, m.Topic)
	case unsubscribe:
		s.Unsubscribe(&client, m.Topic)
	default:
		s.Send(&client, "Server: Action unrecognized")
	}

	return s
}

func (s *PubSub) Publish(sender Client, topic string, message []byte) {
	var clients []Client

	for _, sub := range s.Subscriptions {
		if sub.Topic == topic {
			for _, cli := range *sub.Clients {
				if cli.ID != sender.ID {
					clients = append(clients, cli)
				}
			}
		}
	}

	for _, client := range clients {
		s.Send(&client, string(message))
	}
}

func (s *PubSub) Subscribe(client *Client, topic string) {
	topicExist := false
	clientExist := false

	for _, sub := range s.Subscriptions {
		if sub.Topic == topic {
			topicExist = true
		}

		for _, cli := range *sub.Clients {
			if cli.ID == client.ID {
				clientExist = true
			}
		}
	}

	for _, sub := range s.Subscriptions {
		if sub.Topic == topic && !clientExist {
			topicExist = true
			*sub.Clients = append(*sub.Clients, *client)
		}
	}

	if !topicExist {
		newClient := &[]Client{*client}

		newTopic := &Subscription{
			Topic:   topic,
			Clients: newClient,
		}

		s.Subscriptions = append(s.Subscriptions, *newTopic)
	}
}

func (s *PubSub) Unsubscribe(client *Client, topic string) {
	for _, sub := range s.Subscriptions {
		if sub.Topic == topic {
			for i := 0; i < len(*sub.Clients); i++ {
				if client.ID == (*sub.Clients)[i].ID {
					if i == len(*sub.Clients)-1 {
						*sub.Clients = (*sub.Clients)[:len(*sub.Clients)-1]
					} else {
						*sub.Clients = append((*sub.Clients)[:i], (*sub.Clients)[i+1:]...)
						i--
					}
				}
			}
		}
	}
}
