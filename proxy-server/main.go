package main

import (
	"encoding/json"
	"net/http"
	"os"
	"proxy-server/rabbit_mq"
	"proxy-server/requests"
	"proxy-server/utils"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var mapRequest *requests.Map

func handler(w http.ResponseWriter, req *http.Request) {

	var request requests.Request
	request.Request = req
	request.Response = w
	request.Id = uuid.New().String()

	err := utils.CheckConnection(strings.Split(os.Getenv("url"), "http://")[1])
	if err != nil {
		rabbit_mq.SendRequest("request", request.Id)
		getMap().SetRequest(request.Id, request)
	} else {
		requests.SendRequest(w, req)
	}
}

func ProxyHandler(w http.ResponseWriter, req *http.Request) {
	var request requests.Request
	json.NewDecoder(req.Body).Decode(&request)
	getMap().SendQueuedRequest(request.Id, w, req)
}

func initMap() {
	if mapRequest == nil {
		mapRequest = requests.NewMap()
	}
}

func getMap() *requests.Map {
	return mapRequest
}

func main() {

	utils.ReadFile("./setup.yml")
	initMap()
	router := mux.NewRouter()
	router.HandleFunc("/requests", ProxyHandler)
	router.HandleFunc("/webhook", handler)
	httpHandler := cors.Default().Handler(router)
	http.ListenAndServe(":9091", httpHandler)

}
