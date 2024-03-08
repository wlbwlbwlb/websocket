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

	defer func() {
		if r := recover(); r != nil {
			//log.Error
			fmt.Println(r, debug.Stack(), req)
		}
	}()

	if e = json.Unmarshal(msg, &req); e != nil {
		return
	}

	if f, ok := handlers[req.Cmd]; ok {
		if e = f(req.Payload, c); e != nil {
			return
		}
	}

	return fmt.Errorf("handler not found, req=%+v", req)
}

func init() {
	handlers["ping"] = func(payload []byte, c *Client) (e error) {
		a := struct {
			Cmd string `json:"cmd"`
		}{
			Cmd: "pong",
		}
		resp, _ := json.Marshal(a)
		c.send <- resp
		return
	}
}
