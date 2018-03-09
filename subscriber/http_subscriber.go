package subscriber

import (
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"encoding/json"
	"github.com/3tnet/laravel-broadcasting-server-go/support"
	"gopkg.in/go-playground/validator.v9"
)

type httpSubscriber struct {
	support.Writable
	subCallback BroadcastFunc
	validate    *validator.Validate
}

func (s *httpSubscriber) Subscribe(callback BroadcastFunc) {
	s.subCallback = callback
}

func (s *httpSubscriber) getMessage(data map[string]interface{}) Message {
	var event, socketId string

	if name, ok := data["name"]; ok {
		event = name.(string)
	}

	if sid, ok := data["socket_id"]; ok {
		socketId = sid.(string)
	}
	// var dataStr string
	//if d, ok := data["data"]; ok {
	//	dataStr = d.(string)
	//}

	message := Message{
		Event:    event,
		Data:     data["data"],
		SocketId: socketId,
	}
	return message
}

type Body struct {
	Name     string   `json:"name" validate:"required"`
	Data     string   `json:"data" validate:"required"`
	Channels []string `json:"channels" validate:"required,gte=1"`
	SocketId string   `json:"socket_id"`
}

func (s *httpSubscriber) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if s.subCallback == nil {
		s.Write(rw, r, s.successResponse())
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.Write(rw, r, err)
		return
	}

	body := &Body{}
	if err = json.Unmarshal(b, body); err != nil {
		s.Write(rw, r, err)
		return
	}

	if err = s.validate.Struct(body); err != nil {
		s.Write(rw, r, support.NewHTTPError(http.StatusUnprocessableEntity, err))
		return
	}

	msg := Message{
		Event:    body.Name,
		SocketId: body.SocketId,
	}

	if err = json.Unmarshal([]byte(body.Data), &msg.Data); err != nil {
		s.Write(rw, r, err)
		return
	}

	for _, channel := range body.Channels {
		s.subCallback(channel, msg)
	}

	s.Write(rw, r, s.successResponse())
}

func (s *httpSubscriber) successResponse() interface{} {
	return support.Json(map[string]string{"message": "ok"})
}

func NewHttpSubscriber(router *mux.Router) Subscriber {
	s := &httpSubscriber{validate: validator.New()}
	router.Use(support.WithWriterHandler)
	// 如果 laravel pusher 没配置 appId，这里的 appId 会为空
	router.Handle("/apps/{appId}/events", s).Methods("POST")
	return s
}
