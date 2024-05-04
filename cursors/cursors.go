package cursors

import (
	"net/http"

	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type CursorOS string

const (
	CursorOSMacOS   CursorOS = "macOS"
	CursorOSWindows CursorOS = "windows"
	CursorOSLinux   CursorOS = "linux"
)

type Cursor struct {
	id   string
	os   string
	url  string
	posX int
	posY int
}

// Subscribe subscribes to cursor updates.
//
//encore:api public raw method=GET path=/sub
func Subscribe(w http.ResponseWriter, req *http.Request) {
	id := uuid.New().String()
	ctx := rlog.With("cursor_id", id)

	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		ctx.Error("error upgrading websocket connection", "err", err)
		return
	}
	defer func() {
		ctx.Debug("closing websocket connection")
		c.Close()
	}()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				return
			}
			ctx.Error("error reading message", "err", err)
			return
		}

		if mt != websocket.TextMessage {
			ctx.Error("unexpected message type", "type", mt)
			return
		}

		ctx.Debug("received message", "message", string(message))
	}
}

var upgrader = websocket.Upgrader{}

var db = sqldb.NewDatabase("cursors", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
