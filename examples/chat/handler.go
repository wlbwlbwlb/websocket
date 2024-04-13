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
			fmt.Println(r, string(debug.Stack()), req)
		}
	}()

	f, _ := getHandler(req.Cmd)
	e = f(req.Payload, c)

	return
}

func getHandler(cmd string) (f HandlerFunc, ok bool) {
	if f, ok = handlers[cmd]; ok {
		return
	}
	return nop, true
}

var nop = func(payload []byte, c *Client) (e error) {
	resp := struct {
		Cmd string `json:"cmd"`
	}{
		Cmd: "nop",
	}
	msg, _ := json.Marshal(resp)
	c.send <- msg
	return
}

func init() {
	handlers["ping"] = func(payload []byte, c *Client) (e error) {
		resp := struct {
			Cmd string `json:"cmd"`
		}{
			Cmd: "pong",
		}
		msg, _ := json.Marshal(resp)
		c.send <- msg
		return
	}
}
