package defaults

import "github.com/harmony-one/go-raptorq/pkg/raptorq"
import "github.com/harmony-one/go-raptorq/internal/impl/libraptorq"

func DefaultEncoderFactory() raptorq.EncoderFactory {
	return &libraptorq.EncoderFactory{}
}

func DefaultDecoderFactory() raptorq.DecoderFactory {
	return &libraptorq.DecoderFactory{}
}
