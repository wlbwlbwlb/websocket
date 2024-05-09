package main

var _g *Hub

func Close(reason []byte) (e error) {
	_g.exit <- reason
	return
}
