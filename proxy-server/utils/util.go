package utils

import (
	"log"
	"net"
	"time"
)

func CheckConnection(address string) error {
	timeout := time.Second
	_, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		log.Println("uncreachable : ", err)
		return err
	}
	return nil
}
