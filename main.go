package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/ramseydsilva/snarl"
)

var addr = flag.String("addr", ":8080", "http service address")
var port = flag.Int("port", 9666, "port")
var templ = template.Must(template.ParseFiles("index.html"))

func handleMessage(m snarl.Message) {
	h.broadcast <- []byte(fmt.Sprintf("[%d] %s: %s\n", m.DateTime, string(m.Sender[:]), string(m.Message[:])))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templ.Execute(w, r.Host)
}

func main() {
	go h.run()

	messageChannel := snarl.Receive(*port)
	go func() {
		for {
			handleMessage(<-messageChannel)
		}
	}()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
