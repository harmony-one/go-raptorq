package libraptorq

import (
	"errors"
	"runtime"

	"github.com/harmony-one/go-raptorq/internal/impl/libraptorq/swig"
	"github.com/harmony-one/go-raptorq/internal/readyblockchan"
	"github.com/harmony-one/go-raptorq/pkg/raptorq"
)

// DecoderFactory is a factory of libRaptorQ-based decoder instances.
type DecoderFactory struct {
}

// New returns a new decoder instance.
//
// commonOTI and schemeSpecificOTI are the RaptorQ OTIs,
// received from the sender.
//
// New returns a nil instance and an error if the decoder cannot be created.
// This can, for example,
// occur if the given commonOTI or schemeSpecificOTI is out of range.
func (*DecoderFactory) New(commonOTI uint64, schemeSpecificOTI uint32) (
	decoder raptorq.Decoder, err error) {
	wrapped := swig.NewBytesDecoder(swig.HostToNet64(commonOTI),
		swig.HostToNet32(schemeSpecificOTI))
	if wrapped.Initialized() {
		dec := new(Decoder)
		dec.wrapped = wrapped
		dec.commonOTI = commonOTI
		dec.schemeSpecificOTI = schemeSpecificOTI
		dec.rbcs.Reset(dec.NumSourceBlocks())
		go dec.readyBlocksLoop()
		decoder = dec
		runtime.SetFinalizer(decoder, finalizeDecoder)
	} else {
		swig.DeleteBytesDecoder(wrapped)
		err = errors.New("libRaptorQ decoder failed to initialize")
	}
	return
}

func finalizeDecoder(decoder *Decoder) {
	err := decoder.Close()
	if err != nil {
		// Do nothing for now.  Maybe log in verbose mode once we have one.
	}
}

// Decoder is a RaptorQ decoder instance.
type Decoder struct {
	wrapped           swig.BytesDecoder
	commonOTI         uint64
	schemeSpecificOTI uint32
	rbcs              readyblockchan.ReadyBlockChannels
}

// Decoder destroy sequence:
//
// 1. Decoder loses all references
// 2. GC kicks in
// 3. finalizeDecoder() gets called
// 4. Wrapped decoder is deleted (via Close())
// 5. Decoder dtor waits for thread pool to drain.
// 6. Pending wait_threads() calls see that the object is being deleted and
//    return Error::EXITING via their promise/future pairs.
// 7. Decoder dtor now returns.  Close() and finalizeDecoder() return in turn.
// 8. The future in a WaitForBlock() call (made from readyBlocksLoop()) returns
//    Error::EXITING.
// 9. readyBlocksLoop() sees Error::EXITING and breaks out of loop.
//
// Note that by the time readyBlocksLoop() sees Error::EXITING,
// the “wrapped” field has already been reset as nil.

func (dec *Decoder) readyBlocksLoop() {
	for {
		var sbn uint8
		var e swig.RaptorQ__v1Error
		swig.WaitForBlock(dec.wrapped, &sbn, &e)
		switch e {
		case swig.Error_NONE:
			dec.rbcs.AddBlock(sbn)
		case swig.Error_EXITING:
			break
		}
	}
}

// CommonOTI returns the common object transmission information for the codec.
func (dec *Decoder) CommonOTI() uint64 {
	return dec.commonOTI
}

// TransferLength returns the size of the transfer object, in octets.
func (dec *Decoder) TransferLength() uint64 {
	return dec.commonOTI >> 24
}

// SymbolSize returns the symbol size, in octets.
func (dec *Decoder) SymbolSize() uint16 {
	return uint16(dec.commonOTI)
}

// SchemeSpecificOTI returns the scheme-specific object transmission
// information for the codec.
func (dec *Decoder) SchemeSpecificOTI() uint32 {
	return dec.schemeSpecificOTI
}

// NumSourceBlocks returns the number of source blocks in the transfer object.
func (dec *Decoder) NumSourceBlocks() uint8 {
	return uint8(dec.schemeSpecificOTI >> 24)
}

// SourceBlockSize returns the size of the given source block, in octets,
func (dec *Decoder) SourceBlockSize(sbn uint8) uint32 {
	return uint32(dec.wrapped.Block_size(sbn))
}

// NumSourceSymbols returns the number of source symbols in the given block.
func (dec *Decoder) NumSourceSymbols(sbn uint8) uint16 {
	return dec.wrapped.Symbols(sbn)
}

// NumSubBlocks returns the number of sub-blocks in the given source block.
//
// This is also the same as number of sub-symbols per symbol.
func (dec *Decoder) NumSubBlocks() uint16 {
	return uint16(dec.schemeSpecificOTI >> 8)
}

// SymbolAlignmentParameter returns the symbol alignment parameter, that is,
// the number of octets to which all symbols,
// and sub-symbols should align in memory.
func (dec *Decoder) SymbolAlignmentParameter() uint8 {
	return uint8(dec.schemeSpecificOTI)
}

// Decode decodes the given symbol.
//
// Decoding is done asynchronously,
// so IsSourceObjectReady or IsSourceBlockReady may not immediately return up
// to date result.
func (dec *Decoder) Decode(sbn uint8, esi uint32, symbol []byte) {
	dec.wrapped.Add_symbol(symbol, esi, sbn)
}

// IsSourceBlockReady returns whether the given source block is ready.
func (dec *Decoder) IsSourceBlockReady(sbn uint8) bool {
	return dec.wrapped.Is_block_ready(sbn)
}

// IsSourceObjectReady returns whether the entire source object is ready.
func (dec *Decoder) IsSourceObjectReady() bool {
	return dec.wrapped.Is_ready()
}

// SourceBlock retrieves the given source block into the given buffer.
func (dec *Decoder) SourceBlock(sbn uint8, buf []byte) (n int, err error) {
	n = int(dec.wrapped.Decode_block_bytes(buf, 0, sbn))
	if n != int(dec.SourceBlockSize(sbn)) {
		err = errors.New("decode failure")
	}
	return
}

// SourceObject retrieves the entire source object into the given buffer.
func (dec *Decoder) SourceObject(buf []byte) (n int, err error) {
	n = int(dec.wrapped.Decode_bytes(buf, 0))
	if n != int(dec.TransferLength()) {
		err = errors.New("decode failure")
	}
	return
}

// FreeSourceBlock frees all internal memory used for the given source block.
func (dec *Decoder) FreeSourceBlock(sbn uint8) {
	dec.wrapped.Free(sbn)
}

// AddReadyBlockChan adds a channel through which the decoder notifies the
// block number of each source block ready.
//
// Use this to get notified of source blocks as soon as they become ready.
//
// Source blocks already ready at the time of the call are immediately sent
// to the channel.
//
// AddReadyBlockChan returns an error if the channel has already been added.
func (dec *Decoder) AddReadyBlockChan(ch chan<- uint8) (err error) {
	return dec.rbcs.AddChannel(ch)
}

// RemoveReadyBlockChan removes a channel previously registered using
// AddReadyBlockChan.
//
// RemoveReadyBlockChan returns an error if the channel has not yet been added.
func (dec *Decoder) RemoveReadyBlockChan(ch chan<- uint8) (err error) {
	return dec.rbcs.RemoveChannel(ch)
}

// Close closes the decoder.
func (dec *Decoder) Close() (err error) {
	switch wrapped := dec.wrapped.(type) {
	case swig.BytesDecoder:
		dec.rbcs.Reset(dec.NumSourceBlocks())
		dec.wrapped = nil
		swig.DeleteBytesDecoder(wrapped)
	default:
		err = errors.New("RaptorQ decoder already closed")
	}
	return
}
