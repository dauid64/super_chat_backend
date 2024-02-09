package websockets

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dauid64/super_chat_backend/src/database"
	"github.com/dauid64/super_chat_backend/src/models"

	"github.com/gorilla/websocket"
)

var ChatWebSocket = []WebSocket{
	{
		URI: "/chat",
		Func: serveWsChat,
	},
}

type Client struct {
	Conn  *websocket.Conn
	Email string
}

type MessageWsData struct {
	Type string         `json:"type"`
	User string         `json:"user,omitempty"`
	Message models.Message `json:"message,omitempty"`
}

var clients = make(map[*Client]bool)
var broadcast = make(chan *models.Message)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func receiver(client *Client) {
	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		messageWsData := &MessageWsData{}

		err = json.Unmarshal(p, messageWsData)
		if err != nil {
			log.Println("error while unmarshaling message", err)
			continue
		}

		log.Println("host", client.Conn.RemoteAddr())
		if messageWsData.Type == "bootup" {
			client.Email = messageWsData.User
			log.Println("client successfully mapped", &client, client, client.Email)
		} else {
			log.Println("received message", messageWsData.Type, messageWsData.Message)
			message := messageWsData.Message

			record := database.Instance.Create(&message)
			if record.Error != nil {
				log.Println(record.Error)
				return
			}

			var messageSend models.Message
			record = database.Instance.Joins("ToUser").Joins("FromUser").Where("messages.id = ?", message.ID).Find(&messageSend)
			if record.Error != nil {
				log.Println(record.Error)
				return
			}
			
			broadcast <- &messageSend
		}
	}
}

func broadcaster() {
	for {
		message := <-broadcast
		log.Println("new message", message)

		for client := range clients {
			log.Println("username:", client.Email,
				"fromUser:", message.FromUser,
				"toUser:", message.ToUser)
			
			if (client.Email == message.FromUser.Email || client.Email == message.ToUser.Email) {
				err := client.Conn.WriteJSON(message)
				if err != nil {
					log.Printf("Websocket error: %s", err)
					client.Conn.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func serveWsChat(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := &Client{Conn: ws}

	clients[client] = true
	log.Println("clients", len(clients), clients, ws.RemoteAddr())

	receiver(client)

	log.Println("exiting", ws.RemoteAddr().String())
	delete(clients, client)
}