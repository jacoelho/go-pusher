package main

import (
	"crypto/tls"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	quit    chan struct{}
	conn    *websocket.Conn
	display bool
}

func (c *Client) Read() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		if c.display {
			log.Printf("recv: %s", message)
		} else {
			log.Printf("recv: with len %d", len(message))
		}
	}
}

func (c *Client) Write() {
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case t := <-ticker.C:
			// send a keepalive msg to keep channel open
			err := c.conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}

		// quit channel close, exit
		case _, ok := <-c.quit:
			if !ok {
				log.Println("interrupt")
				// To cleanly close a connection, a client should send a close
				// frame and wait for the server to close the connection.
				err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close:", err)
					return
				}
				select {
				case <-time.After(time.Second):
				}

				ticker.Stop()
				c.conn.Close()
				return
			}
		}
	}

}

func main() {

	numberClients := flag.Int("clients", 10, "number of parallel customers")
	url := flag.String("url", "", "url to connect")
	displayContent := flag.Bool("print", false, "print messages")

	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	quit := make(chan struct{}, 1)
	signal.Notify(interrupt, os.Interrupt)

	d := websocket.Dialer{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	for i := 0; i < *numberClients; i++ {
		c, _, err := d.Dial(*url, nil)
		if err != nil {
			log.Fatal("dial:", err)
		}
		ws := &Client{
			conn:    c,
			quit:    quit,
			display: *displayContent,
		}

		go ws.Read()
		go ws.Write()

	}

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			close(quit)
			select {
			case <-time.After(time.Second * 5):
			}
			return
		}
	}
}
