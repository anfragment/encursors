package cursors

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"encore.dev/metrics"
	"encore.dev/rlog"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type CursorEnterWebsocketPayload struct {
	Id      string   `json:"id"`
	Country string   `json:"country"`
	OS      CursorOS `json:"os"`
	PosX    int      `json:"posX"`
	PosY    int      `json:"posY"`
}

type CursorMoveWebsocketPayload struct {
	Id   string `json:"id"`
	PosX int    `json:"posX"`
	PosY int    `json:"posY"`
}

type CursorLeaveWebsocketPayload struct {
	Id string `json:"id"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var TotalClients = metrics.NewCounter[uint64]("total_clients", metrics.CounterConfig{})

// Subscribe subscribes to cursor updates for a given URL.
//
//encore:api public raw method=GET path=/subscribe
func Subscribe(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	query := req.URL.Query()

	url := query.Get("url")
	if !validateURL(url) {
		http.Error(w, "specify valid url in url parameters", http.StatusBadRequest)
		return
	}

	var os CursorOS
	switch query.Get("os") {
	case "mac":
		os = CursorOSMacOS
	case "win":
		os = CursorOSWindows
	case "linux":
		os = CursorOSLinux
	default:
		http.Error(w, "specify valid os in url parameters", http.StatusBadRequest)
		return
	}

	country := query.Get("country")
	if country == "" || len(country) != 2 {
		http.Error(w, "specify valid country in url parameters", http.StatusBadRequest)
		return
	}

	var posX, posY int
	var err error
	if query.Get("posX") != "" {
		posX, err = strconv.Atoi(query.Get("posX"))
		if err != nil {
			http.Error(w, "specify valid posX in url parameters", http.StatusBadRequest)
			return
		}
	}
	if query.Get("posY") != "" {
		posY, err = strconv.Atoi(query.Get("posY"))
		if err != nil {
			http.Error(w, "specify valid posY in url parameters", http.StatusBadRequest)
			return
		}
	}

	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		rlog.Error("error upgrading websocket connection", "err", err)
		return
	}

	TotalClients.Increment()

	id := uuid.New().String()
	rlog := rlog.With("cursor_id", id)

	events := make(chan *CursorEvent)
	done := make(chan struct{})
	go handleIncomingPubSubEvents(ctx, rlog, id, url, events, done, c)
	defer handleClosure(ctx, rlog, events, done, c)

	cursor := &Cursor{
		Id:      id,
		Country: country,
		OS:      os,
		URL:     url,
		PosX:    posX,
		PosY:    posY,
	}
	if err := writeCursorToDB(ctx, cursor); err != nil {
		rlog.Error("error writing cursor to db", "err", err)
		return
	}
	defer deleteCursorFromDB(ctx, id)

	handleWSComms(ctx, rlog, cursor, c)
}

func handleWSComms(ctx context.Context, rlog rlog.Ctx, cursor *Cursor, c *websocket.Conn) {
	event := &CursorEvent{
		Type:   CursorEventTypeEnter,
		Cursor: cursor,
	}
	if msgId, err := CursorEvents.Publish(ctx, event); err != nil {
		rlog.Error("error publishing cursor event", "err", err)
	} else {
		rlog.Debug("published cursor enter event", "msg_id", msgId)
	}

	var lastMessageAt time.Time

	for {
		if ctx.Err() != nil {
			break
		}

		mt, message, err := c.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				rlog.Error("error reading message", "err", err)
			}
			break
		}
		if mt != websocket.TextMessage {
			rlog.Error("unexpected message type", "type", mt)
			break
		}

		if time.Since(lastMessageAt) < time.Duration(cfg.MinEventTimeoutMs())*time.Millisecond {
			rlog.Error("received message too quickly, ignoring", "elapsed", time.Since(lastMessageAt))
			lastMessageAt = time.Now()
			continue
		}
		lastMessageAt = time.Now()

		pos := [2]int{}
		if err := json.Unmarshal(message, &pos); err != nil {
			rlog.Error("error unmarshalling message", "err", err)
			break
		}

		rlog.Debug("received cursor position", "pos", pos)

		cursor.PosX = pos[0]
		cursor.PosY = pos[1]

		if err := updateCursorInDB(ctx, cursor); err != nil {
			rlog.Error("error updating cursor in db", "err", err)
			break
		}

		event = &CursorEvent{
			Type:   CursorEventTypeMove,
			Cursor: cursor,
		}
		if msgId, err := CursorEvents.Publish(ctx, event); err != nil {
			rlog.Error("error publishing cursor event", "err", err)
		} else {
			rlog.Debug("published cursor move event", "msg_id", msgId)
		}
	}

	event.Type = CursorEventTypeLeave
	if msgId, err := CursorEvents.Publish(ctx, event); err != nil {
		rlog.Error("error publishing cursor event", "err", err)
	} else {
		rlog.Debug("published cursor leave event", "msg_id", msgId)
	}
}

func handleIncomingPubSubEvents(_ context.Context, rlog rlog.Ctx, id string, url string, eventsCh chan *CursorEvent, doneCh <-chan struct{}, c *websocket.Conn) {
	subToUpdates(id, url, eventsCh, doneCh)
	for {
		select {
		case event := <-eventsCh:
			if event.Cursor.Id == id {
				continue
			}
			msg := struct {
				Type    CursorEventType `json:"type"`
				Payload interface{}     `json:"payload"`
			}{
				Type: event.Type,
			}
			switch event.Type {
			case CursorEventTypeEnter:
				msg.Payload = CursorEnterWebsocketPayload{
					Id:      event.Cursor.Id,
					Country: event.Cursor.Country,
					OS:      event.Cursor.OS,
					PosX:    event.Cursor.PosX,
					PosY:    event.Cursor.PosY,
				}
			case CursorEventTypeMove:
				msg.Payload = CursorMoveWebsocketPayload{
					Id:   event.Cursor.Id,
					PosX: event.Cursor.PosX,
					PosY: event.Cursor.PosY,
				}
			case CursorEventTypeLeave:
				msg.Payload = CursorLeaveWebsocketPayload{
					Id: event.Cursor.Id,
				}
			default:
				rlog.Error("unknown cursor event type", "type", event.Type)
				continue
			}

			if err := c.WriteJSON(msg); err != nil {
				rlog.Error("error writing JSON", "err", err)
				return
			}
		case <-doneCh:
			return
		}
	}
}

func handleClosure(ctx context.Context, rlog rlog.Ctx, events chan *CursorEvent, done chan struct{}, c *websocket.Conn) {
	rlog.Debug("closing websocket connection")
	close(done)
	close(events)
	c.Close()
}
