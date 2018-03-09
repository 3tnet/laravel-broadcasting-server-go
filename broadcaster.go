package broadcasting

import (
	"github.com/3tnet/laravel-broadcasting-server-go/subscriber"
	"github.com/ty666/go-socket.io"
)

type Broadcaster struct {
	ioServer *socketio.Server
}

func (b *Broadcaster) Broadcast(channel string, message subscriber.Message) {
	if message.SocketId == "" {
		b.ioServer.BroadcastTo(channel, message.Event, channel, message.Data)
	} else {
		b.ioServer.BroadcastToIgnoreId(message.SocketId, channel, message.Event, channel, message.Data)
	}
}

func (b *Broadcaster) Listen(subs ...subscriber.Subscriber) {
	for _, sub := range subs {
		sub.Subscribe(b.Broadcast)
	}
}

func NewBroadcaster(ioServer *socketio.Server) *Broadcaster {
	return &Broadcaster{ioServer: ioServer}
}
