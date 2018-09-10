package rooms

import (
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

func init() {
	go checker()
}

type Room struct {
	Turn         int
	Game         interfaces.Game
	FirstSocket  *websocket.Conn
	SecondSocket *websocket.Conn
}

var rooms = make(map[string]Room)

var queue struct {
	Body []string
	mux  sync.Mutex
}

/**
	removes room from map and queue(if its exist there)
 */
func closeRoom(token string) {
	r, _ := rooms[token]
	if r.FirstSocket != nil {
		r.FirstSocket.Close()
	}
	if r.SecondSocket != nil {
		r.SecondSocket.Close()
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
func checker() {
	for {
		time.Sleep(time.Second * 5)
		for token, r := range rooms {
			if clientClosedConnection(r.FirstSocket) || clientClosedConnection(r.SecondSocket) {
				closeRoom(token)
			}
		}
	}
}

/**
	pings client and tries to recieve message
 */
func clientClosedConnection(conn *websocket.Conn) bool {
	if conn == nil {
		return false
	}
	conn.WriteJSON("?")
	_, _, err := conn.ReadMessage()
	if err != nil {
		return true
	}
	return false
}

/**
	send turn params to game, switches player and notifies players if game ended
 */
func (r *Room) TurnHandler(params ... int) (bool, error) {
	if params[0] != r.Turn {
		return false, errors.New("not this players turn")
	}
	ok, err := r.Game.Turn(params...)
	if err != nil {
		return false, err
	}
	if r.Turn == 1 {
		r.Turn = 2
	} else {
		r.Turn = 1
	}
	if winner := r.Game.GetWinner(); winner != 0 {
		r.FirstSocket.WriteJSON("W" + strconv.Itoa(winner))
		r.SecondSocket.WriteJSON("W" + strconv.Itoa(winner))
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
				if clientClosedConnection(r.FirstSocket) {
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
	newGame := Room{1, games.GetInstance(gameName), nil, nil}
	newGame.Game.Initialize()
	rooms[token] = newGame
	queue.mux.Lock()
	queue.Body = append(queue.Body, token)
	queue.mux.Unlock()
	return token
}

func GetRoom(token string) (Room, bool) {
	res, ok := rooms[token]
	return res, ok
}

func SetRoom(token string, val Room) {
	rooms[token] = val
}

/**
	links connection to room by token
 */
func NewConnection(token string, conn *websocket.Conn) {
	if r, ok := rooms[string(token)]; ok {
		if r.FirstSocket == nil {
			r.FirstSocket = conn
			rooms[token] = r
		} else {
			r.SecondSocket = conn
			rooms[token] = r
			r.FirstSocket.WriteJSON("1")
			_, msg, _ := r.FirstSocket.ReadMessage()
			r.SecondSocket.WriteJSON(string(msg))
		}
	}
}

func generateToken() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 24)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
