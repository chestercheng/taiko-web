package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"taiko-web/models"
	"time"
)

const consonants = "bcdfghjklmnpqrstvwxyz"

var serverStatus = NewServerStatus()
var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type User struct {
	Ws        *websocket.Conn
	Action    string
	Session   string
	OtherUser *User
}

type ServerStatus struct {
	Users   []*User
	Invites map[string]*User
}

func NewServerStatus() *ServerStatus {
	var s ServerStatus
	s.Invites = make(map[string]*User)
	return &s
}

func (u *User) CreateSession() {
	var seed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	session := make([]byte, 5)
	for i := range session {
		session[i] = consonants[seed.Intn(len(consonants))]
	}

	serverStatus.Invites[string(session)] = u
	u.Action = "invite"
	u.Session = string(session)
}

func (u *User) GameEnd() {
	delete(serverStatus.Invites, u.Session)
	gameend := models.Message{Type: "gameend"}

	if u.OtherUser != nil {
		u.OtherUser.Ws.WriteJSON(&gameend)
		u.OtherUser.Action = "ready"
		u.OtherUser.Session = ""
		u.OtherUser.OtherUser = nil
	}

	u.Ws.WriteJSON(&gameend)
	u.Action = "ready"
	u.Session = ""
	u.OtherUser = nil
}

func MultiplayerHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer ws.Close()

	user := User{Ws: ws, Action: "ready"}
	serverStatus.Users = append(serverStatus.Users, &user)

	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			break
		}

		switch user.Action {
		case "ready":
			// Not playing or waiting
			switch msg.Type {
			case "invite":
				if msg.Value == nil {
					// Session invite link requested
					user.CreateSession()
					user.Ws.WriteJSON(models.Message{"invite", user.Session})
				} else if _, ok := serverStatus.Invites[msg.Value.(string)]; ok {
					// Join a session with the other user
					user.OtherUser = serverStatus.Invites[msg.Value.(string)]
					delete(serverStatus.Invites, msg.Value.(string))

					user.OtherUser.OtherUser = &user
					user.Action = "invite"
					user.Session = msg.Value.(string)

					res := models.Message{Type: "session"}
					user.OtherUser.Ws.WriteJSON(&res)
					user.Ws.WriteJSON(&res)
				} else {
					// Session code is invalid
					user.GameEnd()
				}
			case "join":
				if msg.Value == nil {
					continue
				}

				id := msg.Value.(map[string]interface{})["id"]
				diff := msg.Value.(map[string]interface{})["diff"]
				if id == nil || diff == nil {
					continue
				}

				// Wait for another user
				user.Action = "waiting"
				user.Ws.WriteJSON(models.Message{Type: "waiting"})
			}
		case "invite":
			switch msg.Type {
			case "leave":
				// Cancel session invite
				user.GameEnd()
			case "songsel":
				if user.OtherUser == nil {
					continue
				}

				user.OtherUser.Action = "songsel"
				user.Action = "songsel"

				res := models.Message{Type: "songsel"}
				user.OtherUser.Ws.WriteJSON(&res)
				user.Ws.WriteJSON(&res)
			}
		case "songsel":
			// Session song selection
			if user.OtherUser == nil {
				// Other user disconnected
				user.GameEnd()
				continue
			}

			switch msg.Type {
			case "songsel":
				// Change song select position
				user.OtherUser.Ws.WriteJSON(msg)
				user.Ws.WriteJSON(msg)
			case "join":
				// Start game
				if msg.Value == nil {
					continue
				}

				id := msg.Value.(map[string]interface{})["id"]
				diff := msg.Value.(map[string]interface{})["diff"]
				if id == nil || diff == nil {
					continue
				}

				if user.OtherUser.Action == "waiting" {
					user.OtherUser.Action = "loading"
					user.Action = "loading"

					res := models.Message{"gameload", diff}
					user.OtherUser.Ws.WriteJSON(&res)
					user.Ws.WriteJSON(&res)
				} else {
					user.Action = "waiting"

					var val []map[string]interface{}
					val = append(val, map[string]interface{}{"id": id, "diff": diff})
					user.OtherUser.Ws.WriteJSON(models.Message{"users", val})
				}
			case "gameend":
				// User wants to disconnect
				user.GameEnd()
			}
		case "waiting":
			// Waiting for another user
			switch msg.Type {
			case "leave":
				// Stop waiting
				if user.OtherUser == nil {
					user.GameEnd()
					continue
				}

				if user.Session != "" {
					user.OtherUser.Ws.WriteJSON(models.Message{Type: "users"})
					user.Action = "songsel"
				} else {
					user.Action = "ready"
				}

				(*user.Ws).WriteJSON(models.Message{Type: "left"})
			}
		case "loading":
			switch msg.Type {
			case "gamestart":
				user.OtherUser.Action = "playing"
				user.Action = "playing"

				res := models.Message{Type: "gamestart"}
				user.OtherUser.Ws.WriteJSON(&res)
				user.Ws.WriteJSON(&res)
			}
		case "playing":
			if user.OtherUser == nil {
				// Other user disconnected
				user.GameEnd()
				continue
			}
			switch msg.Type {
			case "note", "drumroll", "branch", "gameresults":
				user.OtherUser.Ws.WriteJSON(msg)
			case "songsel":
				user.OtherUser.Action = "songsel"
				user.Action = "songsel"

				res := models.Message{Type: "songsel"}
				user.OtherUser.Ws.WriteJSON(&res)
				user.Ws.WriteJSON(&res)

				res.Type = "users"
				user.OtherUser.Ws.WriteJSON(&res)
				user.Ws.WriteJSON(&res)
			case "gameend":
				// User wants to disconnect
				user.GameEnd()
			}
		}
	}

	// User disconnected
	user.GameEnd()
}
