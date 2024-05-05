package cursors

import (
	"context"
	"sync"

	"encore.dev/pubsub"
)

type CursorEventType string

const (
	CursorEventTypeEnter CursorEventType = "enter"
	CursorEventTypeMove  CursorEventType = "move"
	CursorEventTypeLeave CursorEventType = "leave"
)

type CursorEvent struct {
	Type   CursorEventType `json:"type"`
	Cursor *Cursor         `json:"cursor"`
}

var subscribersMu sync.RWMutex

// subscribers maps a path to a map of subscriber IDs to channels.
var subscribers = make(map[string]map[string]chan *CursorEvent)

// fanout sends the cursor to all subscribers.
func fanout(ctx context.Context, event *CursorEvent) error {
	subscribersMu.RLock()
	defer subscribersMu.RUnlock()
	for _, ch := range subscribers[event.Cursor.Path] {
		ch <- event
	}
	return nil
}

// subToUpdates subscribes a client to cursor updates.
func subToUpdates(id string, path string, ch chan *CursorEvent, done <-chan struct{}) {
	subscribersMu.Lock()
	defer subscribersMu.Unlock()
	if _, ok := subscribers[path]; !ok {
		subscribers[path] = make(map[string]chan *CursorEvent)
	}
	subscribers[path][id] = ch
	go func() {
		<-done
		subscribersMu.Lock()
		defer subscribersMu.Unlock()
		delete(subscribers[path], id)
	}()
}

var CursorEvents = pubsub.NewTopic[*CursorEvent]("cursor-events", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})

var _ = pubsub.NewSubscription[*CursorEvent](CursorEvents, "cursor-events-fanout", pubsub.SubscriptionConfig[*CursorEvent]{
	Handler: fanout,
})
