package utils

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/go-yaml/yaml"
)

type Setup struct {
	RabbitMQ string `yaml:"rabbitmq"`
	Url      string `yaml:"url"`
}

func ReadFile(file string) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return err
	}
	var info Setup
	if err := yaml.Unmarshal(content, &info); err != nil {
		log.Println(err)
		return err
	}

	os.Setenv("rabbitmq", info.RabbitMQ)
	os.Setenv("url", info.Url)
	return nil
}
