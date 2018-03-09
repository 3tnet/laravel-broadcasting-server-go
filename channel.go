package broadcasting

import (
	"strings"
	"net/http"
	"net/url"
	"github.com/ty666/go-socket.io"
	"github.com/3tnet/laravel-broadcasting-server-go/support"
	"github.com/pkg/errors"
	"log"
)

type Channel struct {
	authHost, authEndpoint string
	privateChannel         privateChannel
	presenceChannel        presenceChannel
	clientEvent            clientEvent
	errLogger              *log.Logger
}

func (c *Channel) auth(channel string, r *http.Request, headers url.Values) support.HTTPError {
	form := url.Values{}
	form.Add("channel_name", channel)
	req, err := http.NewRequest("POST", c.authHost+c.authEndpoint, strings.NewReader(form.Encode()))

	req.Header.Set("Cookie", r.Header.Get("Cookie"))
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	// 必须设置 Content-Type request 中的 body 才会生效
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for key, values := range headers {
		for _, v := range values {
			req.Header.Add(key, v)
		}

	}
	if err != nil {
		return support.NewHTTPError(http.StatusInternalServerError, err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return support.NewHTTPError(http.StatusInternalServerError, err)
	}

	if res.StatusCode != http.StatusOK {
		return support.NewHTTPError(res.StatusCode, errors.New("Client can not be authenticated, got HTTP status:"+string(res.StatusCode)))
	}
	return nil
}

func (c *Channel) joinPrivateAndPresence(channel string, so socketio.Socket, headers url.Values) {
	if err := c.auth(channel, so.Request(), headers); err != nil {
		var statusCode int
		if httpErr, ok := err.(support.HTTPError); ok {
			statusCode = httpErr.StatusCode()
		} else {
			statusCode = http.StatusInternalServerError
		}
		c.errLogger.Println(err)
		so.Emit("subscription_error", channel, statusCode)
	} else {
		so.Join(channel)
		if c.presenceChannel.is(channel) {
			// todo
		}
	}
}

func (c *Channel) Join(channel string, so socketio.Socket, headers url.Values) {
	if channel != "" {
		if c.isPrivateChannel(channel) {
			c.joinPrivateAndPresence(channel, so, headers)
		} else {
			so.Join(channel)
		}
	}
}

func (c *Channel) Leave(channel string, so socketio.Socket) {
	if channel != "" {
		if c.presenceChannel.is(channel) {
			// todo
		}
		so.Leave(channel)
	}
}

func (c *Channel) isPrivateChannel(channel string) bool {
	return c.privateChannel.is(channel) || c.presenceChannel.is(channel)
}

func isInChannel(channel string, so socketio.Socket) bool {
	for _, r := range so.Rooms() {
		if r == channel {
			return true
		}
	}
	return false
}

func (c *Channel) ClientEvent(channel string, so socketio.Socket, event string, data interface{}) {
	if event != "" && channel != "" {
		if c.clientEvent.is(event) && c.isPrivateChannel(channel) && isInChannel(channel, so) {
			so.BroadcastTo(channel, event, channel, data)
		}
	}
}

type privateChannel struct {
}

func (privateChannel) is(channel string) bool {
	return strings.HasPrefix(channel, "private-");
}

type presenceChannel struct {
}

func (presenceChannel) is(channel string) bool {
	return strings.HasPrefix(channel, "presence-");
}

func NewChannel(authHost, authEndpoint string, errLogger *log.Logger) *Channel {
	return &Channel{
		authHost:        authHost,
		authEndpoint:    authEndpoint,
		privateChannel:  privateChannel{},
		presenceChannel: presenceChannel{},
		clientEvent:     clientEvent{},
		errLogger:       errLogger,
	}
}

type clientEvent struct {
}

func (clientEvent) is(event string) bool {
	return strings.HasPrefix(event, "client-")
}
