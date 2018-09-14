package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-rooms/games"
	"go-rooms/lib/rooms"
	"io/ioutil"
	"log"
)

/**
handler.
Returns JSON token and role as struct
*/
func findGame(c *gin.Context) {
	gameName, ok := c.Params.Get("game")
	if !ok {
		c.AbortWithStatus(400)
		return
	}
	c.Writer.Write([]byte(rooms.FindRoom(gameName)))
	c.AbortWithStatus(200)
}

/**
handler.
Accepts parameters as sarray of ints.
Returns http codes
*/
func action(c *gin.Context) {
	token, ok1 := c.Params.Get("room")
	player, ok2 := c.Params.Get("player")
	data, err := c.GetRawData()
	if !(ok1 && ok2) || err != nil {
		c.AbortWithStatus(400)
		return
	}

	r, ok := rooms.GetRoom(token)
	if !ok {
		c.AbortWithStatus(404)
		return
	}
	err = r.Send(player, data)
	if err != nil {
		c.AbortWithStatus(200)
		return
	}
	c.JSON(200, "")
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

/**
Web-socket handler
Handle connection and stores it to specific room
*/
func wshandler(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	roomToken, ok := c.Params.Get("room")
	if !ok {
		c.AbortWithStatus(400)
		return
	}
	playerToken, ok := c.Params.Get("player")
	if ok {
		rooms.NewPingConnection(string(roomToken), string(playerToken), conn)
	} else {
		rooms.NewMessageConnection(string(roomToken), conn)
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// websocket
	router.GET("/ws/:room/", wshandler)
	router.GET("/ws/:room/:player/", wshandler)

	router.GET("/new/:game/", findGame)
	router.POST("/api/:room/action/:player/", action)

	router.StaticFile("/shiii/", "public/template.tmpl")
	router.Static("/assets", "public/assets/")

	// presentation files ----------------------------

	router.LoadHTMLGlob("public/template.tmpl")

	router.GET("/1/", func(c *gin.Context) {
		c.HTML(200, "template.tmpl", gin.H{"content": "tic_tac_toe"})
	})

	router.GET("/2/", func(c *gin.Context) {
		c.HTML(200, "template.tmpl", gin.H{"content": "chatroulette"})
	})

	// ------------------------------------------------

	files, err := ioutil.ReadDir("./games")

	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			if games.GetInstance(f.Name()) == nil {
				fmt.Printf("Warning: Unmapped folder %q in /games/ directory. Please write a case in /games/names.go\n", f.Name())
				continue
			}
			router.StaticFile("/"+f.Name()+"/", "public/"+f.Name()+"/main.html")
			router.Static("/"+f.Name()+"/assets/", "public/"+f.Name())
		}
	}

	router.Run(":8080")
}
