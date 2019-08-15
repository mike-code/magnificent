package main

import (
	"strings"
	"net/http"
	"time"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/websocket"
)

func RunWorker(worker wsWorker) {
	defer func () {
		worker.ws.Close()
		workers.Delete(&worker)
	}()

	welcomeMsg := []byte(fmt.Sprintf("{status: '%s', time: %d}", MAP_STATE[state], time.Now().Unix()))

	if err := worker.ws.WriteMessage(websocket.TextMessage, welcomeMsg); err != nil {
		log.Debugf("Worker %p was unable to send message to the socket. %s", &worker, err.Error())
		return
	}

	for {
		lState   := <- worker.newState
		stateStr := MAP_STATE[lState]
		msg      := []byte(fmt.Sprintf("{status: '%s', time: %d}", stateStr, time.Now().Unix()))

		log.Tracef("Worker %p received state change to %s", &worker, stateStr)

		worker.ws.SetWriteDeadline(time.Now().Add(time.Duration(config.Monitor.Timeout) * time.Millisecond))

		if err := worker.ws.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Warnf("Worker %p was unable to deliver message and will now terminate. %s", &worker, err.Error())
			return
		}
	}
}

func RegisterWorker(ws *websocket.Conn) (wsWorker) {
	worker := wsWorker {
		ws: ws,
		newState: make(chan int),
	}

	workers.Store(&worker, worker)

	if verbose {
		var workersList []string

		workers.Range(func (k interface{}, w interface{}) bool {
			workersList = append(workersList, fmt.Sprintf("%p", k.(*wsWorker)))
			return true
		})

		log.Tracef("Registered workers [%s]", strings.Join(workersList, ","))
	}

	return worker
}

func HandleWsConn(w http.ResponseWriter, req *http.Request) {
	ws, err := (&websocket.Upgrader {
		ReadBufferSize:  256,
		WriteBufferSize: 256,
		CheckOrigin: func (r *http.Request) bool {
			return true // yolo
		},
	}).Upgrade(w, req, nil)

	if err != nil {
		log.Warnf("There was a problem upgrading TCP connection from %s. %s", req.RemoteAddr, err.Error())
		return
	}

	log.Debugf("New WS connection opened from %s", req.RemoteAddr)

	worker := RegisterWorker(ws)

	RunWorker(worker)
}

func HandleSite(w http.ResponseWriter, req *http.Request) {

}

func SpinWebSocket() {
	http.HandleFunc("/",   HandleSite)
	http.HandleFunc("/ws", HandleWsConn)

	s := &http.Server{
		Addr:           config.Monitor.Listen,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}
