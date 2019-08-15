package main

import (
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var workers sync.Map

type wsWorker struct {
	ws *websocket.Conn
	newState chan int
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

		workers.Range(func (k interface{}, w interface{}) (bool) {
			log.Infof("SENDING NEW MSG to %p", k.(*wsWorker))

			w.(wsWorker).newState <- s

			return true
		})
	}
}

func StateListener(stateCheck chan int, stateChange chan int) {
	queue := make([]int, config.Tries.History)

	for {
		lastCheck := <-stateCheck
		newState  := state

		// enqeue new state and dequeue the oldest
		queue = append(queue, lastCheck)[1:]

		if NOW_ALIVE == lastCheck && STATE_UP != state {
			if STATE_TRANS_DOWN == state || config.Tries.Up == sum_islice(queue[len(queue) - config.Tries.Up: ]) {
				newState = STATE_UP
			} else {
				newState = STATE_TRANS_UP
			}
		} else if NOW_DEAD == lastCheck && STATE_DOWN != state {
			if STATE_TRANS_UP == state || 0 == sum_islice(queue[len(queue) - config.Tries.Down: ]) {
				newState = STATE_DOWN
			} else {
				newState = STATE_TRANS_DOWN
			}
		}

		if newState != state {
			log.Infof("State changed [%s -> %s]", MAP_STATE[state], MAP_STATE[newState])

			state = newState
			select {
				// non-blocking channel
				case stateChange <- newState:
				default:
			}
		} else {
			log.Debugf("Current state [%s]", MAP_STATE[state])
		}
	}
}
