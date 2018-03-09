package subscriber

import (
	"testing"
	"github.com/gorilla/mux"
	"net/http/httptest"
	"github.com/gin-gonic/gin/json"
	"bytes"
	"net/http"
)

func TestHttpSubscriber_ServeHTTP(t *testing.T) {
	router := mux.NewRouter()
	hs := NewHttpSubscriber(router)
	hs.Subscribe(func(channel string, message Message) {

	})
	ts := httptest.NewServer(router)
	defer ts.Close()

	body := map[string]interface{}{
		"name":     "App\\Events\\Message",
		"data":     "{\"aa\":\"bb\"}",
		"channels": []string{"chat.1"},
	}
	b, _ := json.Marshal(body)

	r, _ := http.NewRequest("POST", ts.URL+"/apps/1/events", bytes.NewReader(b))

	if res, err := ts.Client().Do(r); err != nil {
		t.Error(err)
	} else {
		if res.StatusCode != http.StatusOK {
			t.Errorf("status code is %d, except %d", res.StatusCode, http.StatusOK)
		}
	}

	body = map[string]interface{}{
		"name":     "App\\Events\\Message",
		"data":     "{\"aa\":\"bb\"}",
		"channels": []string{},
	}
	b, _ = json.Marshal(body)

	r, _ = http.NewRequest("POST", ts.URL+"/apps/1/events", bytes.NewReader(b))

	if res, err := ts.Client().Do(r); err != nil {
		t.Error(err)
	} else {
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("status code is %d, except %d", res.StatusCode, http.StatusUnprocessableEntity)
		}
	}
}
