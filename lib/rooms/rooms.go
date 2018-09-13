package rooms

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"go-rooms/games"
	"go-rooms/lib/interfaces"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Room struct {
	Turn    int
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
		time.Sleep(time.Second / 2)
		if r.Member1.closedConnection() || r.Member2.closedConnection() {
			closeRoom(token)
			break
		}

	}
}

/**
send turn params to game, switches player and notifies players if game ended
*/
func (r *Room) TurnHandler(player string, params ...int) (bool, error) {
	var role = 0
	if r.Member1.Token == player {
		role = 1
	} else if r.Member2.Token == player {
		role = 2
	} else {
		return false, errors.New("at room with token A is no player with token B")
	}
	if role != r.Turn {
		return false, errors.New("not this player's turn")
	}
	ok, err := r.Game.Turn(append([]int{role}, params...)...)
	if err != nil {
		return false, err
	}
	if r.Turn == 1 {
		r.Turn = 2
	} else {
		r.Turn = 1
	}
	if ok {
		msg, _ := json.Marshal(r.Game.GetGrid())
		r.Member1.send(string(msg))
		r.Member2.send(string(msg))
	}
	if winner := r.Game.GetWinner(); winner != 0 {
		r.Member1.send("W" + strconv.Itoa(winner))
		r.Member2.send("W" + strconv.Itoa(winner))
	}
	return ok, nil
}

/**
finds place in room in queue or returns new room
returns game interface and role
*/
func FindRoom(gameName string) (string, int) {
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
				return token, 2
			}
		}
		queue.mux.Unlock()
	}
	return createRoom(gameName), 1
}

/**
creates new room push it to queue
returns token
*/
func createRoom(gameName string) (token string) {
	token = generateToken()
	newGame := Room{1, games.GetInstance(gameName), Client{Alive: false}, Client{Alive: false}}
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
		if !r.Member1.Alive {
			r.Member1.msgConn = conn
			r.Member1.Token = generateToken()
			r.Member1.send(r.Member1.Token)
		} else {
			if r.Member2.msgConn != nil {
				r.Member2.pingConn = conn
				r.Member2.Alive = true
				r.Member1.send("Connected")
				r.Member2.send("Connected")
			} else {
				r.Member2.msgConn = conn
				r.Member2.Token = generateToken()
				r.Member2.send(r.Member2.Token)
			}
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
			r.Member1.send("Connected")
			r.Member2.send("Connected")
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
