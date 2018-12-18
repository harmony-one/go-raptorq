package libraptorq

import (
	"errors"
	"runtime"

	"github.com/harmony-one/go-raptorq/internal/impl/libraptorq/swig"
	"github.com/harmony-one/go-raptorq/pkg/raptorq"
)

type EncoderFactory struct {
}

func (*EncoderFactory) New(input []byte, symbolSize uint16, minSubSymbolSize uint16,
	maxSubBlockSize uint32, alignment uint8) (enc raptorq.Encoder, err error) {
	wrapped := swig.InitBytesEncoder(input, minSubSymbolSize, symbolSize,
		int64(maxSubBlockSize))
	if !wrapped.Initialized() {
		swig.DeleteBytesEncoder(wrapped)
		err = errors.New("libRaptorQ encoder failed to initialize")
	} else {
		enc = &Encoder{wrapped, maxSubBlockSize}
		runtime.SetFinalizer(enc, finalizeEncoder)
	}
	return
}

func finalizeEncoder(encoder *Encoder) {
	encoder.Close()
}

type Encoder struct {
	wrapped         swig.BytesEncoder
	maxSubBlockSize uint32
}

func (enc *Encoder) CommonOTI() uint64 {
	return swig.NetToHost64(enc.wrapped.OTI_Common())
}

func (enc *Encoder) TransferLength() uint64 {
	return enc.CommonOTI() >> 24
}

func (enc *Encoder) SymbolSize() uint16 {
	return uint16(enc.CommonOTI())
}

func (enc *Encoder) SchemeSpecificOTI() uint32 {
	return uint32(swig.NetToHost32(enc.wrapped.OTI_Scheme_Specific()))
}

func (enc *Encoder) NumSourceBlocks() uint8 {
	return uint8(enc.SchemeSpecificOTI() >> 24)
}

func (enc *Encoder) SourceBlockSize(sbn uint8) uint32 {
	return uint32(enc.wrapped.Block_size(sbn))
}

func (enc *Encoder) NumSourceSymbols(sbn uint8) uint16 {
	return enc.wrapped.Symbols(sbn)
}

func (enc *Encoder) NumSubBlocks() uint16 {
	return uint16(enc.SchemeSpecificOTI() >> 8)
}

func (enc *Encoder) SymbolAlignmentParameter() uint8 {
	return uint8(enc.SchemeSpecificOTI())
}

func (enc *Encoder) Encode(sbn uint8, esi uint32, buf []byte) (written uint, err error) {
	if len(buf) < int(enc.SymbolSize()) {
		err = errors.New("RaptorQ encoder buffer too small")
	} else {
		written = uint(enc.wrapped.Encode(buf, esi, sbn))
		if written == 0 {
			err = errors.New("RaptorQ encoder returned an error indication")
		}
	}
	return
}

func (enc *Encoder) MaxSubBlockSize() uint32 {
	return enc.maxSubBlockSize
}

func (enc *Encoder) FreeSourceBlock(sbn uint8) {
	enc.wrapped.Free(sbn)
}

func (enc *Encoder) MinSymbols(sbn uint8) uint16 {
	return uint16(enc.wrapped.Extended_symbols(sbn))
}

func (enc *Encoder) MaxSymbols(sbn uint8) uint32 {
	return uint32(enc.wrapped.Max_repair(sbn))
}

func (enc *Encoder) Close() (err error) {
	switch wrapped := enc.wrapped.(type) {
	case swig.BytesEncoder:
		swig.DeleteBytesEncoder(wrapped)
		enc.wrapped = nil
	default:
		err = errors.New("RaptorQ encoder already closed")
	}
	return
}
