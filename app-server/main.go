package main

import (
	"app-server/rabbit_mq"
	"app-server/utils"
	"bytes"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, req *http.Request) {

	var b = &bytes.Buffer{}
	err := req.Write(b)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(b)
}

func main() {

	utils.ReadFile("./setup.yml")

	go func() {
		rabbit_mq.Receive()
	}()

	http.HandleFunc("/", handler)
	http.ListenAndServe(":9080", nil)

}
