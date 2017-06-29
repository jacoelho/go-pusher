package service

import (
	"encoding/json"
	"log"

	"bitbucket.org/jacoelho/pusher/message"
)

func Validate(in <-chan []byte) <-chan []byte {
	out := make(chan []byte)

	go func() {
		for {
			select {
			case m := <-in:
				var msg message.Platforms

				if err := json.Unmarshal(m, &msg); err != nil {
					log.Println(err)
					break
				}

				// message ok, forward
				out <- m
			default:

			}
		}
	}()

	return out
}
