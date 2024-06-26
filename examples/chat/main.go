// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"flag"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wlbwlbwlb/log"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	hub := newHub()
	go hub.run()

	g = hub

	//if e := mq.Init(mq.Lookupd("localhost:4161"),
	//	mq.Nsqd("127.0.0.1:4150"),
	//); e != nil {
	//	log.Fatal(e)
	//}
	//defer func() {
	//	mq.StopConsumers()
	//	mq.StopProducer()
	//}()

	//http.HandleFunc("/", serveHome)
	//http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	//	serveWs(hub, w, r)
	//})
	r := gin.Default()
	r.GET("/", serveHome)
	r.GET("/ws", func(c *gin.Context) {
		serveWs(g, c)
	})
	server := &http.Server{
		Addr:    *addr,
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.L.Fatalf("listen: %s", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.L.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.L.Fatalf("server shutdown: %s", err)
	}

	//这种停服方式有问题 todo
	//hub.shutdown()

	log.L.Info("server exiting")
}

//func serveHome(w http.ResponseWriter, r *http.Request) {
//	log.Println(r.URL)
//	if r.URL.Path != "/" {
//		http.Error(w, "Not found", http.StatusNotFound)
//		return
//	}
//	if r.Method != http.MethodGet {
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//	http.ServeFile(w, r, "home.html")
//}

func serveHome(c *gin.Context) {
	homeTemplate.Execute(c.Writer, "ws://"+c.Request.Host+"/ws")
}

var homeTemplate, _ = loadHtml()

func loadHtml() (tmpl *template.Template, e error) {
	t, e := os.ReadFile("examples/chat/home.html")
	if e != nil {
		panic(e)
	}
	tmpl = template.Must(template.New("").Parse(string(t)))
	return
}
