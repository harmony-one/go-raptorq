package libraptorq

import (
	"errors"
	"github.com/harmony-one/go-raptorq/pkg/impl/libraptorq/swig"
	"github.com/harmony-one/go-raptorq/pkg/raptorq"
	"runtime"
)

type DecoderFactory struct {
}

func (*DecoderFactory) New(commonOTI uint64, schemeSpecificOTI uint32) (
	decoder raptorq.Decoder, err error) {
	wrapped := swig.NewBytesDecoder(swig.HostToNet64(commonOTI),
		swig.HostToNet32(schemeSpecificOTI))
	if wrapped.Initialized() {
		dec := new(Decoder)
		dec.wrapped = wrapped
		dec.commonOTI = commonOTI
		dec.schemeSpecificOTI = schemeSpecificOTI
		decoder = dec
		runtime.SetFinalizer(decoder, finalizeDecoder)
	} else {
		swig.DeleteBytesDecoder(wrapped)
		err = errors.New("libRaptorQ decoder failed to initialize")
	}
	return
}

func finalizeDecoder(decoder *Decoder) {
	decoder.Close()
}

type Decoder struct {
	wrapped           swig.BytesDecoder
	commonOTI         uint64
	schemeSpecificOTI uint32
}

func (dec *Decoder) CommonOTI() uint64 {
	return dec.commonOTI
}

func (dec *Decoder) TransferLength() uint64 {
	return dec.commonOTI >> 24
}

func (dec *Decoder) SymbolSize() uint16 {
	return uint16(dec.commonOTI)
}

func (dec *Decoder) SchemeSpecificOTI() uint32 {
	return dec.schemeSpecificOTI
}

func (dec *Decoder) NumSourceBlocks() uint8 {
	return uint8(dec.schemeSpecificOTI >> 24)
}

func (dec *Decoder) SourceBlockSize(sbn uint8) uint32 {
	return uint32(dec.wrapped.Block_size(sbn))
}

func (dec *Decoder) NumSourceSymbols(sbn uint8) uint16 {
	return dec.wrapped.Symbols(sbn)
}

func (dec *Decoder) NumSubBlocks() uint16 {
	return uint16(dec.schemeSpecificOTI >> 8)
}

func (dec *Decoder) SymbolAlignmentParameter() uint8 {
	return uint8(dec.schemeSpecificOTI)
}

func (dec *Decoder) Decode(sbn uint8, esi uint32, symbol []byte) {
	dec.wrapped.Add_symbol(symbol, esi, sbn)
}

func (dec *Decoder) IsSourceBlockReady(sbn uint8) bool {
	return dec.wrapped.Is_block_ready(sbn)
}

func (dec *Decoder) IsSourceObjectReady() bool {
	return dec.wrapped.Is_ready()
}

func (dec *Decoder) SourceBlock(sbn uint8, buf []byte) (n int, err error) {
	n = int(dec.wrapped.Decode_block_bytes(buf, 0, sbn))
	if n != int(dec.SourceBlockSize(sbn)) {
		err = errors.New("decode failure")
	}
	return
}

func (dec *Decoder) SourceObject(buf []byte) (n int, err error) {
	n = int(dec.wrapped.Decode_bytes(buf, 0))
	if n != int(dec.TransferLength()) {
		err = errors.New("decode failure")
	}
	return
}

func (dec *Decoder) FreeSourceBlock(sbn uint8) {
	dec.wrapped.Free(sbn)
}

func (dec *Decoder) Close() (err error) {
	switch wrapped := dec.wrapped.(type) {
	case swig.BytesDecoder:
		dec.wrapped = nil
		swig.DeleteBytesDecoder(wrapped)
	default:
		err = errors.New("RaptorQ decoder already closed")
	}
	return
}
