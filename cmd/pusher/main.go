package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jacoelho/pusher/emiter/redis"
	"github.com/jacoelho/pusher/handlers"
	"github.com/jacoelho/pusher/service"
	"github.com/jacoelho/pusher/websocket"
)

func main() {
	redisUrl := flag.String("redis-server", "127.0.0.1:6379", "redis url")
	redisPort := flag.Int("redis-database", 0, "redis port")
	redisChannel := flag.String("redis-channel", "test", "redis channel")

	flag.Parse()

	logger := log.New(os.Stdout, "server: ", log.Lshortfile)

	r := redis.New(*redisUrl, *redisPort, logger)

	if err := r.Watch(*redisChannel); err != nil {
		logger.Fatal(err)
	}

	hub := websocket.NewHub(logger)

	// start hub main loop
	go hub.RunWithEmiter(service.Validate(r.Events()))

	http.Handle("/ws", handlers.WithLogger(logger, http.HandlerFunc(hub.ServeWebsocket)))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Fatal("ListenAndServe: ", err)
	}
}
