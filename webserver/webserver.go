package webserver

import (
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/leaanthony/mewn"
	"harbored/config"
	"harbored/database"
	"harbored/models/presentation"
	"net/http"
	"time"
)

var Server *echo.Echo

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func init() {
	// Get client files as binary
	clientJs := mewn.MustBytes("./client/dist/main.js")
	clientHtml := mewn.MustBytes("./client/dist/index.html")
	// Room initialization
	globalRoom := NewRoom("global")
	go globalRoom.init()
	RoomList.Set("global", globalRoom)
	// Server start
	Server = echo.New()
	Server.HideBanner = true
	// Middleware
	Server.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	Server.Use(middleware.Logger())
	Server.Use(middleware.Recover())
	Server.Use(sentryecho.New(sentryecho.Options{
		Repanic: true,
	}))
	// Routing
	Server.Static("/static", config.Config.StaticDir)
	Server.GET("/api/presentations", func(c echo.Context) error {
		txn := database.DB.Txn(false)
		defer txn.Abort()
		iterator, err := txn.Get("presentations", "id")
		if err != nil {
			return err
		}
		var presentations []presentation.Presentation
		for obj := iterator.Next(); obj != nil; obj = iterator.Next() {
			p := obj.(presentation.Presentation)
			presentations = append(presentations, p)
		}
		return c.JSON(http.StatusOK, presentations)
	})
	Server.GET("/", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "text/html", clientHtml)
	})
	Server.GET("/main.js", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/javascript", clientJs)
	})
	Server.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	Server.GET("/api/ws", wsEntry)
	//go Server.Start(config.Config.ServerPort)
}

// Setup WebSocket connection
func wsEntry(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	client := NewClient(ws)
	defer client.close()
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	go client.read()
	for {
		select {
		// Send message to a client
		case message := <-client.Send:
			client.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			err := client.Connection.WriteMessage(1, message)
			if err != nil {
				return nil
			}
		// Ping a client
		case <-ticker.C:
			client.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			err := client.Connection.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				return nil
			}
		}
	}
}
