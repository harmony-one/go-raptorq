// Package readyblockchan provides a mix-in that implements ready-block
// channel interface of raptorq.Decoder.
package readyblockchan

import (
	"fmt"
	"sync"
)

type AlreadyAdded chan<- uint8

func (e AlreadyAdded) Error() string {
	return fmt.Sprintf("ready-block channel %+v already added")
}

type NotFound chan<- uint8

func (e NotFound) Error() string {
	return fmt.Sprintf("ready-block channel %+v not found")
}

// ReadyBlockChannels is a collection of ready-block channels.
type ReadyBlockChannels struct {
	mutex    sync.Mutex
	ready    []bool
	channels []chan<- uint8
}

func (rbcs *ReadyBlockChannels) Reset(numSourceBlocks uint8) {
	rbcs.mutex.Lock()
	defer rbcs.mutex.Unlock()
	for _, ch := range rbcs.channels {
		close(ch)
	}
	rbcs.channels = nil
	rbcs.ready = make([]bool, numSourceBlocks)
}

func (rbcs *ReadyBlockChannels) AddChannel(ch chan<- uint8) (
	err error,
) {
	rbcs.mutex.Lock()
	defer rbcs.mutex.Unlock()
	for _, ch1 := range rbcs.channels {
		if ch == ch1 {
			err = AlreadyAdded(ch)
			return
		}
	}
	rbcs.channels = append(rbcs.channels, ch)
	for sbn, ready := range rbcs.ready {
		if ready {
			go addBlockToChan(uint8(sbn), ch)
		}
	}
	return
}

func (rbcs *ReadyBlockChannels) RemoveChannel(ch chan<- uint8) (
	err error,
) {
	rbcs.mutex.Lock()
	defer rbcs.mutex.Unlock()
	var idx = -1
	for i, ch1 := range rbcs.channels {
		if ch == ch1 {
			idx = i
			break
		}
	}
	if idx == -1 {
		err = NotFound(ch)
		return
	}
	rbcs.channels[idx] = rbcs.channels[len(rbcs.channels)-1]
	rbcs.channels = rbcs.channels[:len(rbcs.channels)-1]
	return
}

func (rbcs *ReadyBlockChannels) AddBlock(sbn uint8) {
	rbcs.mutex.Lock()
	defer rbcs.mutex.Unlock()
	if rbcs.ready[sbn] {
		return
	}
	rbcs.ready[sbn] = true
	for _, ch := range rbcs.channels {
        go addBlockToChan(sbn, ch)
	}
}

func addBlockToChan(sbn uint8, ch chan<- uint8) {
    ch <- sbn
}
