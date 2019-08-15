package main

import (
	"os/signal"
	"os"
	"errors"
	"flag"
	"io/ioutil"
	"bufio"
	"bytes"
	"strconv"
	"net/http"
	"time"
	"net"
	"fmt"
	log "github.com/sirupsen/logrus"
)

var (
	config  Config
	debug   = false
	verbose = false

	MAP_STATE = map[int]string {
		STATE_UP:         "UP",
		STATE_DOWN:       "DOWN",
		STATE_TRANS_UP:   "T_UP",
		STATE_TRANS_DOWN: "T_DOWN",
	}

	MAP_CHECK = map[int]string {
		NOW_DEAD:  "ERR",
		NOW_ALIVE: "OK",
	}
)

const (
	STATE_UP         = 0
	STATE_DOWN       = 1
	STATE_TRANS_UP   = 2
	STATE_TRANS_DOWN = 3

	NOW_DEAD         = 0
	NOW_ALIVE        = 1
)

func ValidateHttpResponse(buffer []byte) (error) {
	response, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(buffer)), nil)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not parse HTTP response. %s", err.Error()))
	}

	state.httpStatus = response.StatusCode

	if config.Http.Validate.Status != response.StatusCode {
		return errors.New(fmt.Sprintf("Wrong status code. Expect %d got %d", config.Http.Validate.Status, response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not read HTTP response body. %s", err.Error()))
	}

	if !config.Http.Validate.regex.Match(body) {
		return errors.New("HTTP Response body doesn't match given regex.")
	}

	return nil
}

func HttpCheck(conn net.Conn) (error) {
	var buf = make([]byte, config.Chunk)

	reqBody := fmt.Sprintf("%s %s HTTP/%s\r\n", config.Http.Method, config.Http.Query, config.Http.Version)

	if config.Http.Version != "1.0" {
		reqBody += fmt.Sprintf("Host: %s\r\n", config.Hostname)
	}

	reqBody += "\r\n"
	log.Tracef(">> Request\n%s", reqBody)

	conn.SetDeadline(time.Now().Add(time.Duration(config.Timeout.Check) * time.Millisecond))

	if _, err := conn.Write([]byte(reqBody)); err != nil {
		return err
	}

	if _, err := conn.Read(buf); err != nil {
		return err
	}

	log.Tracef("<< Response\n%s", string(buf))

	if config.Http.Validate.Enabled {
		if err := ValidateHttpResponse(buf); err != nil {
			return err
		}
	}

	return nil
}

func TcpConn() (net.Conn, error) {
	addr := net.JoinHostPort(config.Hostname, strconv.Itoa(config.Port))
	t_o  := time.Duration(config.Timeout.Connect) * time.Millisecond

	if _, err := net.ResolveIPAddr("ip", config.Hostname); err != nil {
		log.Fatalf("Failed to resolve hostname [%s]. The application will not recover from this error.", config.Hostname)
	}

	log.Tracef("Establishing TCP connection to %s. Timeout set to %dms.", addr, t_o.Nanoseconds()/1e6)

	return net.DialTimeout("tcp", addr, t_o)
}

func DoCheck() (error) {
	t := time.Now()

	defer func() {
		state.checkDuration = time.Now().Sub(t).Nanoseconds()/1e6
	}()

	conn, err := TcpConn()

	if err != nil {
		return errors.New(fmt.Sprintf("Error while establishing a TCP connection. %s", err.Error()))
	}

	defer conn.Close()

	if config.Tcponly == false {
		if err := HttpCheck(conn); err != nil {
			return errors.New(fmt.Sprintf("Error during HTTP poll. %s", err.Error()))
		}
	}

	return nil
}

func RunChecks(stateCheck chan int) {
	for {
		if err := DoCheck(); err == nil {
			log.Trace("Alive")

			stateCheck <- NOW_ALIVE
		} else {
			log.Tracef("Dead: %s", err)

			stateCheck <- NOW_DEAD
		}

		time.Sleep(time.Duration(config.Interval) * time.Millisecond)
	}
}

func main() {
	flag.BoolVar(&debug,   "v",  false, "Enable debugging")
	flag.BoolVar(&verbose, "vv", false, "Enable verbose debugging")
	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	if verbose {
		debug = true
		log.SetLevel(log.TraceLevel)
	}

	log.Info("Started at " + time.Now().String())

	stateCheck  := make(chan int)
	stateChange := make(chan int)

	LoadConfig()

	// initially we assume the website is down
	state.currentState    = STATE_DOWN
	state.lastCheckStatus = NOW_DEAD

	if config.Tcponly {
		state.checkType = "L4"
	} else {
		state.checkType = "L7"
	}

	if config.Monitor.Enabled {
		go SpinWebSocket()
	}

	go StateListener(stateCheck, stateChange)
	go StatePooler(stateChange)
	go RunChecks(stateCheck)

	signalCh := make(chan os.Signal, 1)
    signal.Notify(signalCh, os.Interrupt)
    <-signalCh
}
