package ingestor

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sammyoina/seriesOfTubes/queue"
)

type WebsocketListener struct {
	router *gin.Engine
	route  string
}

func NewWebsocketListener(route string) *WebsocketListener {
	r := gin.Default()
	authConfig := cors.DefaultConfig()
	authConfig.AllowAllOrigins = true
	r.Use(cors.New(authConfig))
	return &WebsocketListener{
		router: r,
		route:  route,
	}
}

func (l *WebsocketListener) StartAccepting(q queue.Queue) {
	l.router.GET(l.route, func(c *gin.Context) {
		var upgrader = websocket.Upgrader{}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			q.Enqueue(message)
		}
	})
	l.router.Run()
}
