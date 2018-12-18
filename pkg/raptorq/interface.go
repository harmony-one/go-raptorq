// raptorq package provides the RaptorQ encoder/decoder and their factory
// interfaces.

package raptorq

// ObjectInfo provides various codec information about the source object.
type ObjectInfo interface {
	// CommonOTI returns the Common FEC Object Transmission Information.
	CommonOTI() uint64

	// TransferLength returns the source object size, in octets.  “F” in RFC
	// 6330.
	TransferLength() uint64

	// SymbolSize returns the symbol size, in octets.  “T” in RFC 6330.
	SymbolSize() uint16

	// SchemeSpecificOTI returns the RaptorQ Scheme-Specific FEC Object
	// Transmission Information.
	SchemeSpecificOTI() uint32

	// NumSourceBlocks returns the number of source blocks.  “Z” in RFC 6330.
	NumSourceBlocks() uint8

	// SourceBlockSize returns the size of the given source block, in bytes, or
	// 0 if sbn is out of range.
	SourceBlockSize(sbn uint8) uint32

	// NumSourceSymbols returns the number of encoding symbols covering the
	// original source data in the given source block, or 0 if sbn is out of
	// range.  The first encoding symbols up to this number is called source
	// symbols, and contains the source block data itself.
	NumSourceSymbols(sbn uint8) uint16

	// NumSubBlocks returns the number of sub-blocks.  “N” in RFC 6330.
	NumSubBlocks() uint16

	// SymbolAlignmentParameter returns the symbol alignment, in octets.  “Al”
	// in RFC 6330.
	SymbolAlignmentParameter() uint8
}

// Encoder encodes one object into a series of symbols.
type Encoder interface {
	// Encoder needs to provide object information.
	ObjectInfo

	// Encode writes the encoding symbol identified by the given source block
	// number and encoding symbol ID, into the given buffer.
	//
	// On success, Encode returns the number of octets written into buf and nil
	// error.
	//
	// On error, Encode returns a non-nil error code.
	Encode(sbn uint8, esi uint32, buf []byte) (written uint, err error)

	// MaxSubBlockSize returns the maximum size block that is decodable in
	// working memory, in octets.  “WS” in RFC 6330.
	MaxSubBlockSize() uint32

	// FreeSourceBlock, on supported implementations, will free memory used for
	// generating encoding symbols for the given source block.  Once a source
	// block has been freed, calling Encode with its SBN may return an error.
	FreeSourceBlock(sbn uint8)

	// MinSymbols returns the minimum number of encoding symbols needed to
	// achieve the 99% probability of decoding the given source block,
	// or 0 if sbn is out of range.  “K” in RFC 6330.
	MinSymbols(sbn uint8) uint16

	// MaxSymbols returns the maximum number of encoding symbols the encoder can
	// generate for the given source block, or 0 if sbn is out of range.
	MaxSymbols(sbn uint8) uint32

	// Close closes the Encoder.  After an Encoder is closed, all methods but
	// Close() will panic if called.
	Close() error
}

// EncoderFactory is a factory of Encoder instances.
type EncoderFactory interface {
	/*
		New creates and returns an Encoder that can encode the given source
		object into symbols.

		input is the source object to encode.

		symbolSize is the encoding symbol size, in octets.

		minSubSymbolSize is the minimum encoding symbol size allowed, in octets.
		(If you are not sure, or if you do not want internal interleaving of
		source symbols, set it equal to symbolSize.)

		maxSubBlockSize is the maximum size block that is decodable in working
		memory, in octets.  Iff this is lower than the source object size, the
		source object will be split into more than one source blocks.  The
		maximum allowed value is 56403 * symbolSize.

		alignment is an internal alignment parameter, in bytes.  Typically this
		is a power of 2, up to implementation-defined maximum.  Larger alignment
		will speed up calculation, at the expense of slightly higher
		transmission size overhead.  Both symbolSize and minSubSymbolSize must
		be a multiple of this.

		On success, New returns an Encoder instance and nil error; on failure,
		it returns nil Encoder and an error code.
	*/
	New(input []byte, symbolSize uint16, minSubSymbolSize uint16,
		maxSubBlockSize uint32, alignment uint8) (Encoder, error)
}

// Decoder decodes encoding symbols and reconstructs one object from a series of
// symbols.
type Decoder interface {
	// Decoder needs to provide object information.
	ObjectInfo

	// Decode decodes a received encoding symbol.
	Decode(sbn uint8, esi uint32, symbol []byte)

	// IsSourceBlockReady returns whether the given source block has been fully
	// decoded and ready to be retrieved, or false if sbn is out of range.
	IsSourceBlockReady(sbn uint8) bool

	// IsSourceObjectReady returns whether the entire source object has been
	// fully decoded and ready to be retrieved.
	IsSourceObjectReady() bool

	// SourceBlock copies the given source block into the given buffer.  buf
	// should contain enough space to store the given source block (use
	// SourceBlockSize(sbn) to get the required size).
	SourceBlock(sbn uint8, buf []byte) (n int, err error)

	// SourceObject copies the source object into the given buffer.  buf should
	// contain enough space to store the entire source object (use
	// TransferLength() to get the required size).
	SourceObject(buf []byte) (n int, err error)

	// Free, on supported implementations, will free memory used for generating
	// encoding symbols for the given source block.  Once a source block has
	// been freed, calling Encode with its SBN may return an error.
	FreeSourceBlock(sbn uint8)

	// Close closes the Decoder.  After a Decoder is closed, all methods but
	// Close() will panic if called.
	Close() error
}

// DecoderFactory is a factory of Decoder instances.
type DecoderFactory interface {
	/*
		New creates and returns a Decoder that can decode incoming source
		symbols and recover the original source object.

		commonOTI is the Common FEC Object Transmission Information, received
		from the sender's encoder.

		schemeSpecificOTI is the Scheme-Specific FEC Object Transmission
		Information, received from the sender's encoder.

		On success, New returns a Decoder instance and nil error; on failure,
		it returns nil Encoder and an error code.
	*/
	New(commonOTI uint64, schemeSpecificOTI uint32) (Decoder, error)
}
