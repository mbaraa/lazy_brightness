package main

import (
	"log"
	"net"
	"net/http"
)

type Router struct {
	handler *http.ServeMux
}

func NewRouter(bcAPI *BCWebAPI) *Router {
	return (&Router{http.NewServeMux()}).initEndPoints(bcAPI)
}

func (r *Router) initEndPoints(bcApi *BCWebAPI) *Router {
	r.handler.Handle("/", http.FileServer(http.Dir(".")))
	r.handler.Handle("/brits/", bcApi)
	return r
}

func (r *Router) Start() {
	log.Printf("running server on http://%s", r.getMachineIP())
	log.Fatalln(
		http.ListenAndServe(":80", r.handler),
	)
}

// Get preferred outbound ip of this machine
func (r *Router) getMachineIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
