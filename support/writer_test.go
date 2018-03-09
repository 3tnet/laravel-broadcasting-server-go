package support

import (
	"encoding/json"
	"errors"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testStruct struct {
	Name string `json:"name"`
}

func (t testStruct) ToJson() ([]byte, error) {
	return json.Marshal(t)
}

func TestWriter_SetContent(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	w := NewWriter()
	w.SetContent("abc1").Write(responseWriter)
	assert.Equal(t, responseWriter.Body.String(), "abc1")
	assert.Equal(t, responseWriter.Code, http.StatusOK)

	responseWriter = httptest.NewRecorder()
	w = NewWriter()
	w.SetContent(123).Write(responseWriter)
	assert.Equal(t, responseWriter.Body.String(), "123")
	assert.Equal(t, responseWriter.Code, http.StatusOK)

	responseWriter = httptest.NewRecorder()
	w = NewWriter()
	w.SetContent([]byte("abc2")).Write(responseWriter)
	assert.Equal(t, responseWriter.Body.String(), "abc2")
	assert.Equal(t, responseWriter.Code, http.StatusOK)

	responseWriter = httptest.NewRecorder()
	w = NewWriter()
	w.SetContent(testStruct{"哦哦哦"}).Write(responseWriter)
	assert.Equal(t, responseWriter.Body.String(), `{"name":"哦哦哦"}`)
	assert.Equal(t, responseWriter.Header().Get("Content-Type"), "application/json")
	assert.Equal(t, responseWriter.Code, http.StatusOK)

	responseWriter = httptest.NewRecorder()
	w = NewWriter()
	w.SetContent(errors.New("abc3")).Write(responseWriter)
	b, _ := NewAPIError(http.StatusInternalServerError, "abc3").ToJson()
	assert.Equal(t, responseWriter.Body.String(), string(b))
	assert.Equal(t, responseWriter.Code, http.StatusInternalServerError)

	responseWriter = httptest.NewRecorder()
	w = NewWriter()
	w.SetContent(NewAPIError(http.StatusBadRequest, "bad")).Write(responseWriter)
	assert.Equal(t, responseWriter.Body.String(), `{"message":"bad"}`)
	assert.Equal(t, responseWriter.Header().Get("Content-Type"), "application/json")
	assert.Equal(t, responseWriter.Code, http.StatusBadRequest)

}

func TestJson(t *testing.T) {
	one := struct {
		Name string
		Age  int
	}{"ty", 21}
	responseWriter := httptest.NewRecorder()
	w := NewWriter()
	w.Write(responseWriter, Json(one))
	assert.Equal(t, responseWriter.Body.String(), `{"Name":"ty","Age":21}`)
	assert.Equal(t, responseWriter.Header().Get("Content-Type"), "application/json")
	assert.Equal(t, responseWriter.Code, http.StatusOK)
}
