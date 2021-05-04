package requests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func SendRequestProxy(id string) {

	payload := map[string]string{"id": id}
	json, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", os.Getenv("url"), bytes.NewBuffer(json))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

}
