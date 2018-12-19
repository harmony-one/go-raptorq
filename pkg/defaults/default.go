package defaults

import "github.com/harmony-one/go-raptorq/pkg/raptorq"
import "github.com/harmony-one/go-raptorq/internal/impl/libraptorq"

// DefaultEncoderFactory is the default encoder factory.
func DefaultEncoderFactory() raptorq.EncoderFactory {
	return &libraptorq.EncoderFactory{}
}

// DefaultDecoderFactory is the default decoder factory.
func DefaultDecoderFactory() raptorq.DecoderFactory {
	return &libraptorq.DecoderFactory{}
}

// NewEncoder creates and returns an encoder using the default factory.
func NewEncoder(
	input []byte, symbolSize uint16, minSubSymbolSize uint16,
	maxSubBlockSize uint32, alignment uint8,
) (enc raptorq.Encoder, err error) {
	factory := DefaultEncoderFactory()
	return factory.New(
		input, symbolSize, minSubSymbolSize, maxSubBlockSize, alignment,
	)
}

// NewDecoder creates and returns a decoder using the default factory.
func NewDecoder(commonOTI uint64, schemeSpecificOTI uint32) (
	dec raptorq.Decoder, err error,
) {
	factory := DefaultDecoderFactory()
	return factory.New(commonOTI, schemeSpecificOTI)
}
