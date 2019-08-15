package main

import (
	"strconv"
	"net"
	"io/ioutil"
	"encoding/json"
	"strings"
	"net/http"
	"time"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/gorilla/websocket"
)

type responseObject struct {
	Time            int64  `json:"time"`
	Type            string `json:"type"`
	LastChkDuration int64  `json:"lastChkDuration"`
	State           string `json:"state"`
	HttpStatus      int    `json:"httpStatus"`
	LastCheck       string `json:"lastCheck"`
}

func CraftWsMsg() ([]byte) {
	res := responseObject {
		Time:            time.Now().Unix(),
		Type:            state.checkType,
		LastChkDuration: state.checkDuration,
		State:           MAP_STATE[state.currentState],
		HttpStatus:      state.httpStatus,
		LastCheck:       MAP_CHECK[state.lastCheckStatus],
	}

	msg, _ := json.Marshal(res)

	return msg
}

func RunWorker(worker *wsWorker) {
	defer func () {
		log.Tracef("Deleting %p", worker)
		workers.Delete(worker)
		worker.ws.Close()
	}()

	if err := worker.ws.WriteMessage(websocket.TextMessage, CraftWsMsg()); err != nil {
		log.Debugf("Worker %p was unable to send message to the socket. %s", worker, err.Error())
		return
	}

	log.Tracef("Starting worker %p", worker)

	for {
		<- worker.newState

		stateStr := MAP_STATE[state.currentState]

		log.Tracef("Worker %p received state change to %s", worker, stateStr)

		worker.ws.SetWriteDeadline(time.Now().Add(time.Duration(config.Monitor.Timeout) * time.Millisecond))

		if err := worker.ws.WriteMessage(websocket.TextMessage, CraftWsMsg()); err != nil {
			log.Warnf("Worker %p was unable to deliver message and will now terminate. %s", worker, err.Error())
			return
		}
	}
}

func RegisterWorker(ws *websocket.Conn) (*wsWorker) {
	worker := wsWorker {
		ws: ws,
		newState: make(chan int),
	}

	log.Tracef("Registering new worker %p", &worker)

	workers.Store(&worker, worker)

	if verbose {
		var workersList []string

		workers.Range(func (k interface{}, w interface{}) bool {
			workersList = append(workersList, fmt.Sprintf("%p", k.(*wsWorker)))
			return true
		})

		log.Tracef("Registered workers [%s]", strings.Join(workersList, ","))
	}

	return &worker
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

func HandleSite(writer http.ResponseWriter, req *http.Request) {
	html, err := ioutil.ReadFile("monitor.html")
	if err != nil {
		log.Fatalf("Could not open monitor.html file. %s", err.Error())
	}

	port := strconv.Itoa(req.Context().Value(http.LocalAddrContextKey).(*net.TCPAddr).Port)

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(writer, strings.Replace(string(html), "__PORT__", port, 1))
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
