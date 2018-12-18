package raptorq

import "github.com/harmony-one/go-raptorq/pkg/defaults"

// NewEncoder creates and returns an encoder using the default factory.
func NewEncoder(
	input []byte, symbolSize uint16, minSubSymbolSize uint16,
	maxSubBlockSize uint32, alignment uint8,
) (enc Encoder, err error) {
	factory := defaults.DefaultEncoderFactory()
	return factory.New(
		input, symbolSize, minSubSymbolSize, maxSubBlockSize, alignment,
	)
}

// NewDecoder creates and returns a decoder using the default factory.
func NewDecoder(commonOTI uint64, schemeSpecificOTI uint32) (
	dec Decoder, err error,
) {
	factory := defaults.DefaultDecoderFactory()
	return factory.New(commonOTI, schemeSpecificOTI)
}
