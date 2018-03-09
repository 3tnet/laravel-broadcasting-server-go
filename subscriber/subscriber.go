package subscriber


type Message struct {
	Event    string
	Data     interface{}
	SocketId string
}

type BroadcastFunc func(channel string, message Message)

type Subscriber interface {
	Subscribe(callback BroadcastFunc)
}
