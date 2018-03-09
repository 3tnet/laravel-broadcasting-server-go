package broadcasting

import (
	"github.com/ty666/go-socket.io"
	"github.com/3tnet/laravel-broadcasting-server-go/subscriber"
	"net/url"
	"log"
)

type Server struct {
	broadcaster *Broadcaster
	channel     *Channel
	logger      *log.Logger
}

func (s *Server) Listen(subs ...subscriber.Subscriber) {
	s.broadcaster.Listen(subs...)
}

func (s *Server) OnConnection(so socketio.Socket) {
	s.onSubscribe(so)
	s.onUnsubscribe(so)
	s.onDisconnecting(so)
	s.onClientEvent(so)
}

func parseSubscribeHeaders(data map[string]interface{}) (headers url.Values) {
	headers = url.Values{}
	if a, ok := data["auth"]; ok {
		if auth, ok := a.(map[string]interface{}); ok {
			if h, ok := auth["headers"]; ok {
				if h1, ok := h.(map[string]interface{}); ok {
					for k, v := range h1 {
						if str, ok := v.(string); ok {
							headers.Add(k, str)
						}
					}
				}
			}
		}
	}
	return
}
func getString(key string, data map[string]interface{}) string {
	if v, ok := data[key]; ok {
		if val, ok := v.(string); ok {
			return val
		}
	}
	return ""
}

func (s *Server) onSubscribe(so socketio.Socket) {
	so.On("subscribe", func(data map[string]interface{}) {
		channel := getString("channel", data)
		s.logger.Printf("subscribe %s channel\n", channel)
		headers := parseSubscribeHeaders(data)
		s.channel.Join(channel, so, headers)
	})
}

func (s *Server) onUnsubscribe(so socketio.Socket) {
	so.On("unsubscribe", func(data map[string]interface{}) {
		channel := getString("channel", data)
		s.logger.Printf("unsubscribe %s channel\n", channel)
		s.channel.Leave(channel, so)
	})
}

func (s *Server) onDisconnecting(so socketio.Socket) {

	// todo 为什么不触发
	so.On("disconnecting", func() {
		for _, channel := range so.Rooms() {
			if channel != so.Id() {
				s.channel.Leave(channel, so)
			}
		}
	})
}

func (s *Server) onClientEvent(so socketio.Socket) {
	so.On("client event", func(data map[string]interface{}) {
		channel := getString("channel", data)
		event := getString("event", data)
		s.logger.Printf("client event, channel: %s, event: %v\n", channel, data)
		s.channel.ClientEvent(channel, so, event, data["data"])
	})
}

func NewServer(ioServer *socketio.Server, authHost, authEndpoint string, logger, errLogger *log.Logger) *Server {
	s := &Server{
		broadcaster: NewBroadcaster(ioServer),
		channel:     NewChannel(authHost, authEndpoint, errLogger),
		logger:      logger,
	}

	ioServer.On("connection", func(so socketio.Socket) {
		s.OnConnection(so)
	})
	return s
}
