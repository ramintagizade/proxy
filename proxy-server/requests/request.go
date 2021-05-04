package requests

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Request struct {
	Id       string              `json:"id"`
	Request  *http.Request       `json:"request"`
	Response http.ResponseWriter `json:"response"`
}

type Map struct {
	RequestsMap map[string]Request
}

func NewMap() *Map {
	m := &Map{}
	m.Init()
	return m
}

func (m *Map) Init() {
	m.RequestsMap = make(map[string]Request)
}

func (m *Map) SetRequest(id string, request Request) {

	new_request := request.Request

	var b bytes.Buffer
	b.ReadFrom(request.Request.Body)
	request.Request.Body = ioutil.NopCloser(&b)
	new_request.Body = ioutil.NopCloser(bytes.NewReader(b.Bytes()))
	m.RequestsMap[id] = Request{
		Id:       id,
		Request:  new_request,
		Response: request.Response,
	}
}

func (m *Map) GetRequest(id string) Request {
	return m.RequestsMap[id]
}

func SendRequest(w http.ResponseWriter, req *http.Request) {
	url, err := url.Parse(os.Getenv("url"))
	if err != nil {
		log.Println("url parse err : ", err)
		return
	}
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host
	httputil.NewSingleHostReverseProxy(url).ServeHTTP(w, req)
}

func (m *Map) SendQueuedRequest(id string, w http.ResponseWriter, r *http.Request) {

	req := m.GetRequest(id)
	url, err := url.Parse(os.Getenv("url"))
	if err != nil {
		log.Println("url parse err : ", err)
		return
	}

	req.Request.URL.Host = url.Host
	req.Request.URL.Scheme = url.Scheme
	req.Request.Header.Set("X-Forwarded-Host", req.Request.Header.Get("Host"))
	req.Request.Host = url.Host

	r.Body = req.Request.Body
	r.ContentLength = req.Request.ContentLength
	r.Header = req.Request.Header
	r.Method = req.Request.Method
	r.Response = req.Request.Response

	serveProxy(url, w, r)
}

func serveProxy(url *url.URL, res http.ResponseWriter, req *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(res, req)
}
