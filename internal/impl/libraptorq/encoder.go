package libraptorq

import (
	"errors"
	"runtime"

	"github.com/harmony-one/go-raptorq/internal/impl/libraptorq/swig"
	"github.com/harmony-one/go-raptorq/pkg/raptorq"
)

// EncoderFactory is a factory of libRaptorQ-based encoder instances.
type EncoderFactory struct {
}

// New creates a new encoder instance.
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
	err := encoder.Close()
	if err != nil {
		// Do nothing for now.  Maybe log in verbose mode once we have one.
	}
}

// Encoder is a libRaptorQ-based encoder instance.
type Encoder struct {
	wrapped         swig.BytesEncoder
	maxSubBlockSize uint32
}

// CommonOTI returns the common object transmission information for the codec.
func (enc *Encoder) CommonOTI() uint64 {
	return swig.NetToHost64(enc.wrapped.OTI_Common())
}

// TransferLength returns the length of the source object, in octets.
func (enc *Encoder) TransferLength() uint64 {
	return enc.CommonOTI() >> 24
}

// SymbolSize returns the size of each symbol, in octets.
func (enc *Encoder) SymbolSize() uint16 {
	return uint16(enc.CommonOTI())
}

// SchemeSpecificOTI returns the scheme-specific object transmission
// information for the codec.
func (enc *Encoder) SchemeSpecificOTI() uint32 {
	return uint32(swig.NetToHost32(enc.wrapped.OTI_Scheme_Specific()))
}

// NumSourceBlocks returns the number of source blocks in the source object.
func (enc *Encoder) NumSourceBlocks() uint8 {
	return uint8(enc.SchemeSpecificOTI() >> 24)
}

// SourceBlockSize returns the size of the given source block, in octets.
func (enc *Encoder) SourceBlockSize(sbn uint8) uint32 {
	return uint32(enc.wrapped.Block_size(sbn))
}

// NumSourceSymbols returns the number of source symbols in the given block.
func (enc *Encoder) NumSourceSymbols(sbn uint8) uint16 {
	return enc.wrapped.Symbols(sbn)
}

// NumSubBlocks returns the number of sub-blocks in the given source block.
func (enc *Encoder) NumSubBlocks() uint16 {
	return uint16(enc.SchemeSpecificOTI() >> 8)
}

// SymbolAlignmentParameter returns the number of octets to which all symbols
// and sub-symbols align in memory.
func (enc *Encoder) SymbolAlignmentParameter() uint8 {
	return uint8(enc.SchemeSpecificOTI())
}

// Encode retrieves one encoding symbol,
// identified by the given source block number – encoding symbol ID pair.
//
// Encode returns the number of octets written into the given buffer,
// and an error indication, or nil if no error.
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

// MaxSubBlockSize returns the maximum sub-block size, in octets.
//
// This number is WS * Al in RFC 6330.
func (enc *Encoder) MaxSubBlockSize() uint32 {
	return enc.maxSubBlockSize
}

// FreeSourceBlock frees resource used for encoding the given source block.
func (enc *Encoder) FreeSourceBlock(sbn uint8) {
	enc.wrapped.Free(sbn)
}

// MinSymbols is the number of encoding symbols that needs to be generated and
// sent for the given source block,
// so that the receiver can retrieve the source block with 99% probability.
//
// This number is K in RFC 6330.
func (enc *Encoder) MinSymbols(sbn uint8) uint16 {
	return uint16(enc.wrapped.Extended_symbols(sbn))
}

// MaxSymbols is the number of encoding symbols that can potentially be
// generated for the given source block.  It is somewhere around 2**24.
func (enc *Encoder) MaxSymbols(sbn uint8) uint32 {
	return uint32(enc.wrapped.Max_repair(sbn))
}

// Close closes the encoder instance.
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
