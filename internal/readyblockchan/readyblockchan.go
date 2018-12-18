// Package readyblockchan provides a mix-in that implements ready-block
// channel interface of raptorq.Decoder.
package readyblockchan

import (
	"fmt"
	"sync"
)

// AlreadyAdded signals the given ready block channel has already been added.
type AlreadyAdded chan<- uint8

func (e AlreadyAdded) Error() string {
	return fmt.Sprintf("ready-block channel %+v already added", e)
}

// NotFound signals the given ready block channel is not found.
type NotFound chan<- uint8

func (e NotFound) Error() string {
	return fmt.Sprintf("ready-block channel %+v not found", e)
}

// ReadyBlockChannels is a collection of ready-block channels.
type ReadyBlockChannels struct {
	mutex    sync.Mutex
	ready    []bool
	channels []chan<- uint8
}

// Reset resets this instance.  Existing channels are closed and removed,
// and all blocks are reset as not received.
func (rbcs *ReadyBlockChannels) Reset(numSourceBlocks uint8) {
	rbcs.mutex.Lock()
	defer rbcs.mutex.Unlock()
	for _, ch := range rbcs.channels {
		close(ch)
	}
	rbcs.channels = nil
	rbcs.ready = make([]bool, numSourceBlocks)
}

// AddChannel adds the given channel.
//
// If any source block has already been received,
// AddChannel sends its number immediately.
//
// If the channel already exists, AddChannel returns an error.
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

// RemoveChannel removes the given channel.
//
// Removed channel is not closed.
//
// If the given channel is not found, RemoveChannel returns an error.
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

// AddBlock adds a source block as having been received.
//
// For each source block,
// the first call – and only the first call – sends the block number to all
// channels registered.
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
