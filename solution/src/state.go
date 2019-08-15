package main

import (
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	workers sync.Map
	state   appState
)

type wsWorker struct {
	ws *websocket.Conn
	newState chan int
}

type appState struct {
	checkType       string
	checkDuration   int64
	httpStatus      int
	lastCheckStatus int
	currentState    int
}

func sum_islice(slice []int) (int) {
	sum := 0

	for _, v := range slice {
		sum += v
	}

	return sum
}

func StatePooler(sChange chan int) {
	for {
		s := <-sChange

		log.Debugf("Trigger status change")

		if verbose {
			cnt := 0
			workers.Range(func (k interface{}, w interface{}) (bool) {
				cnt += 1
				return true
			})
			log.Infof("Workers %d", cnt)
		}

		workers.Range(func (k interface{}, w interface{}) (bool) {
			log.Debugf("Updating state channel for worker %p", k.(*wsWorker))

			w.(wsWorker).newState <- s

			return true
		})
	}
}

func StateListener(stateCheck chan int, stateChange chan int) {
	queue := make([]int, config.Tries.History)

	for {
		lastCheck := <-stateCheck
		newState  := state.currentState

		state.lastCheckStatus = lastCheck

		// enqeue new state and dequeue the oldest
		queue = append(queue, lastCheck)[1:]

		if NOW_ALIVE == lastCheck && STATE_UP != state.currentState {
			if STATE_TRANS_DOWN == state.currentState || config.Tries.Up == sum_islice(queue[len(queue) - config.Tries.Up: ]) {
				newState = STATE_UP
			} else {
				newState = STATE_TRANS_UP
			}
		} else if NOW_DEAD == lastCheck && STATE_DOWN != state.currentState {
			if STATE_TRANS_UP == state.currentState || 0 == sum_islice(queue[len(queue) - config.Tries.Down: ]) {
				newState = STATE_DOWN
			} else {
				newState = STATE_TRANS_DOWN
			}
		}

		if newState != state.currentState {
			log.Infof("State changed [%s -> %s]", MAP_STATE[state.currentState], MAP_STATE[newState])

			state.currentState = newState

			select {
				// non-blocking channel
				case stateChange <- newState:
				default:
			}
		} else {
			log.Debugf("Current state [%s]", MAP_STATE[state.currentState])

			if config.Monitor.TransitionOnly == false {
				select {
					// non-blocking channel
					case stateChange <- state.currentState:
					default:
				}
			}
		}
	}
}
