package main

var h *Hub

func Close(reason []byte) (e error) {
	h.exit <- reason
	return
}
