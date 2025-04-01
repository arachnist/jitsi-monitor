package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type apiServer struct {
	h       *http.Server
	members map[string][]string
	lock    sync.RWMutex
}

func (a *apiServer) getMembers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		a.lock.Lock()
		defer a.lock.Unlock()
		j, _ := json.Marshal(a.members)
		w.Write(j)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "sorry, can't do that")
	}
}

func (a *apiServer) Update(channel string, members []string) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.members[channel] = members
}

func StartAPIServer() *apiServer {
	httpServer := &http.Server{
		Addr: listen,
	}

	members := make(map[string][]string)

	apiServer := &apiServer{
		h:       httpServer,
		members: members,
	}

	http.HandleFunc("/jitsi", apiServer.getMembers)

	go apiServer.h.ListenAndServe()

	return apiServer
}
