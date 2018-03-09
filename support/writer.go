package support

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"errors"
)

type Response interface {
	StatusCode() int
	SetStatusCode(int)
	Header() http.Header
	SetHeader(http.Header)
	Content() []byte
	SetContent([]byte)
}

type Jsonable interface {
	ToJson() ([]byte, error)
}

type response struct {
	statusCode int
	header     http.Header
	content    []byte
}

func (r *response) StatusCode() int {
	return r.statusCode
}
func (r *response) SetStatusCode(statusCode int) {
	r.statusCode = statusCode
}
func (r *response) Header() http.Header {
	return r.header
}
func (r *response) SetHeader(header http.Header) {
	r.header = header
}

func (r *response) Content() []byte {
	return r.content
}

func (r *response) SetContent(content []byte) {
	r.content = content
}

type Writer struct {
	response   Response
	wasWritten bool
}

func (w *Writer) SetResponse(r Response) {
	w.response = r
	w.SetContent(r.Content())
}

func (w *Writer) SetHeader(key, value string) *Writer {
	w.response.Header().Set(key, value)
	return w
}

func (w *Writer) AddHeader(key, value string) *Writer {
	w.response.Header().Add(key, value)
	return w
}

func (w *Writer) SetStatusCode(statusCode int) *Writer {
	w.response.SetStatusCode(statusCode)
	return w
}

func (w *Writer) handleError(e error) []byte {
	httpErr, ok := e.(HTTPError)

	if !ok {
		httpErr = NewAPIError(http.StatusInternalServerError, e.Error())
	}

	w.SetStatusCode(httpErr.StatusCode())

	if _, ok := httpErr.(Jsonable); ok {
		w.SetHeader("Content-Type", "application/json")
	}

	return []byte(httpErr.Error())
}

func (w *Writer) SetContent(content interface{}) *Writer {
	var bytes []byte
	switch content.(type) {
	case error:
		bytes = w.handleError(content.(error))
	case Jsonable:
		bytes, _ = content.(Jsonable).ToJson()
		w.SetHeader("Content-Type", "application/json")
	case []byte:
		bytes = content.([]byte)
	case string:
		bytes = []byte(content.(string))
	default:
		bytes = []byte(fmt.Sprint(content))
	}
	w.response.SetContent(bytes)
	return w
}

func (w *Writer) writeHeader(rw http.ResponseWriter) {
	for k, v := range w.response.Header() {
		for _, val := range v {
			rw.Header().Add(k, val)
		}
	}
}

func (w *Writer) Write(rw http.ResponseWriter, content ...interface{}) (n int, err error) {
	if w.wasWritten {
		return 0, errors.New("Multiple Write")
	}
	w.wasWritten = true

	if len(content) > 0 {
		w.SetContent(content[0])
	}
	// rw.WriteHeader 必须放到 w.writeHeader 下面
	w.writeHeader(rw)
	rw.WriteHeader(w.response.StatusCode())
	return rw.Write(w.response.Content())
}

func NewWriter() *Writer {
	return &Writer{&response{statusCode: http.StatusOK, header: http.Header{}}, false}
}

type writeContextKey struct{}

func FromWriterContext(ctx context.Context) *Writer {
	wt := ctx.Value(writeContextKey{})
	if w, ok := wt.(*Writer); ok {
		return w
	}
	return nil
}

// NewWriterContext returns a new Context carrying Writer.
func NewWriterContext(parent context.Context, w *Writer) context.Context {
	return context.WithValue(parent, writeContextKey{}, w)
}

type jsonWarp struct {
	data interface{}
}

func (j jsonWarp) ToJson() ([]byte, error) {
	if data, ok := j.data.(Jsonable); ok {
		return data.ToJson()
	}
	return json.Marshal(j.data)
}

func Json(data interface{}) Jsonable {
	return jsonWarp{data}
}

func WithWriterHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if FromWriterContext(r.Context()) == nil {
			r = r.WithContext(NewWriterContext(r.Context(), NewWriter()))
		}
		handler.ServeHTTP(w, r)
	})
}

type Writable struct{}

func (w *Writable) Writer(r *http.Request) *Writer {
	wt := FromWriterContext(r.Context())
	if wt == nil {
		panic(errors.New("writer not set"))
	}
	return wt
}

func (w *Writable) Write(rw http.ResponseWriter, r *http.Request, data interface{}) (int, error) {
	return w.Writer(r).SetContent(data).Write(rw)
}
