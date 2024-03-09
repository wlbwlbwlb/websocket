package main

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
)

type Msg struct {
	Cmd     string          `json:"cmd"`
	Payload json.RawMessage `json:"payload"`
}

type HandlerFunc func(payload []byte, c *Client) error

var handlers = make(map[string]HandlerFunc)

func handle(msg []byte, c *Client) (e error) {
	var req Msg
	if e = json.Unmarshal(msg, &req); e != nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			//log.Error
			fmt.Println(r, string(debug.Stack()), req)
		}
	}()

	f, ok := handlers[req.Cmd]
	if !ok {
		return fmt.Errorf("handler not found, req=%+v", req)
	}
	e = f(req.Payload, c)

	return
}

func init() {
	handlers["ping"] = func(payload []byte, c *Client) (e error) {
		a := struct {
			Cmd string `json:"cmd"`
		}{
			Cmd: "pong",
		}
		//panic("aaa")
		resp, _ := json.Marshal(a)
		c.send <- resp
		return
	}
}
