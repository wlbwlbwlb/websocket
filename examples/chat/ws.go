package main

var g *Hub

func Close(reason []byte) (e error) {
	g.exit <- reason
	return
}
