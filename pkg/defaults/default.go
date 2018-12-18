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
