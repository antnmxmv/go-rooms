package rooms

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"go-rooms/games"
	"go-rooms/lib/interfaces"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type Room struct {
	Game    interfaces.Game
	Member1 Client
	Member2 Client
}

var rooms = make(map[string]*Room)

var queue struct {
	Body []string
	mux  sync.Mutex
}

/**
removes room from map and queue(if its exist there)
*/
func closeRoom(token string) {
	r, _ := rooms[token]
	if r.Member1.Alive {
		r.Member1.close()
	}
	if r.Member2.Alive {
		r.Member2.close()
	}
	delete(rooms, token)
	queue.mux.Lock()
	for i, val := range queue.Body {
		if val == token {
			queue.Body = append(queue.Body[:i], queue.Body[i+1:]...)
			break
		}
	}
	queue.mux.Unlock()
}

/*
	daemon. removes rooms, where client closed connection
*/
func checker(token string) {
	r, _ := rooms[token]
	for {
		time.Sleep(time.Second * 2)
		if r.Member1.closedConnection() || r.Member2.closedConnection() {
			closeRoom(token)
			break
		}

	}
}

/**
send turn params to game, switches player and notifies players if game ended
*/
func (r *Room) Send(player string, message json.RawMessage) error {
	if !r.Member2.Alive || !r.Member1.Alive {
		return errors.New("not all players connected")
	}
	var role = 0
	if r.Member1.Token == player {
		role = 1
	} else if r.Member2.Token == player {
		role = 2
	} else {
		return errors.New("at room with token A is no player with token B")
	}
	msg, err := r.Game.Action(role, message)
	if err != nil {
		return err
	}

	r.Member1.send(msg)
	r.Member2.send(msg)
	return nil
}

/**
finds place in room in queue or returns new room
returns game interface and role
*/
func FindRoom(gameName string) string {
	if len(queue.Body) != 0 {
		queue.mux.Lock()
		for i, token := range queue.Body {
			r, _ := rooms[token]
			if strings.Compare(r.Game.GetName(), gameName) == 0 {
				if r.Member1.closedConnection() {
					queue.Body = append(queue.Body[:i], queue.Body[i+1:]...)
					queue.mux.Unlock()
					return FindRoom(gameName)
				}
				queue.Body = append(queue.Body[:i], queue.Body[i+1:]...)
				queue.mux.Unlock()
				return token
			}
		}
		queue.mux.Unlock()
	}
	return createRoom(gameName)
}

/**
creates new room push it to queue
returns token
*/
func createRoom(gameName string) (token string) {
	token = generateToken()
	newGame := Room{games.GetInstance(gameName), Client{Alive: false}, Client{Alive: false}}
	newGame.Game.Initialize()
	rooms[token] = &newGame
	queue.mux.Lock()
	queue.Body = append(queue.Body, token)
	queue.mux.Unlock()
	return token
}

func GetRoom(token string) (*Room, bool) {
	res, ok := rooms[token]
	return res, ok
}

/**
links COMMON connection to member of specified room
*/
func NewMessageConnection(token string, conn *websocket.Conn) {
	if r, ok := rooms[token]; ok {
		if r.Member1.msgConn == nil {
			r.Member1.msgConn = conn
			r.Member1.Token = generateToken()
			r.Member1.send(r.Member1.Token)
		} else {
			r.Member2.msgConn = conn
			r.Member2.Token = generateToken()
			r.Member2.send(r.Member2.Token)
		}
	}
}

/**
links PING/PONG connection to member of specified room
*/
func NewPingConnection(roomToken, playerToken string, conn *websocket.Conn) {
	if r, ok := rooms[roomToken]; ok {
		if r.Member1.Token == playerToken {
			r.Member1.pingConn = conn
			r.Member1.Alive = true
			go checker(roomToken)
		} else if r.Member2.Token == playerToken {
			r.Member2.pingConn = conn
			r.Member2.Alive = true
			r.Member1.send("1")
			r.Member2.send("2")
		}
	}
}

func generateToken() string {
	rand.Seed(time.Now().UTC().UnixNano())
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 24)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
