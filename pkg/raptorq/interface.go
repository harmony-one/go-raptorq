/*
raptorq package provides the RaptorQ encoder/decoder and their factory
interfaces, as well as a few concrete factory implementations.

In order to send a binary object, termed “source object,” the sender creates an
Encoder for the source object.  The Encoder first splits the source object into
one or more contiguous blocks, termed “source blocks,” then for each source
block, the Encoder generates encoded and serially numbered chunks of data,
termed “encoding symbols.”  For each source block, the sender uses the Encoder
to generate as many encoding symbols as needed by the receiver to recover the
source block.  The sender chooses the size of the encoding symbol, in octets,
when creating an Encoder.

On the other side, the receiver creates a Decoder for the same source object,
then keeps feeding the Decoder with the encoding symbols received from the
sender, until the Decoder is able to recover the source object from the
encoding symbols.  Once the source object has been recovered, the receiver
closes and discards the Decoder.

An encoding symbol is a unit data of transmission.  Each encoding symbol is
identified by a (SBN, ESI) pair, where SBN is the 8-bit serial number of the
source block from which the symbol was generated, and ESI is the 24-bit serial
identifier of the encoding symbol within the source block.  Both SBN and ESI are
zero-based.

Upon receipt of an encoding symbol, the receiver first needs to identify the
Decoder to use for the source object to which the received encoding symbol
belongs, then feed the Decoder with the encoding symbol, along with its SBN and
ESI.

In order to recover each source block with high probability, the receiver needs
as many encoding symbols as needed to fulfill the original size (in octets) of
the source block.  For example, a 1MiB source block with 1KiB encoding symbols,
the receiver would need 1024 encoding symbols with 99% probability.  Each
additional symbol adds roughly “two nines” to the decoding success probability:
In the example above, 1025 encoding symbols would mean 99.99%, and 1026 encoding
symbols would mean 99.9999%.

It is okay for an encoding symbol to be completely lost (erased) during
transit.  For each encoding symbol lost, the encoder simply needs to generate
and send another encoding symbol.  The sender need not know which encoding
symbol was lost; a brand new encoding symbol would be able to replace for any
previously sent-and-lost encoding symbol.  The sender may even anticipate losses
and send additional encoding symbols in advance without having to wait for
negative acknowledgements (NAKs) from the receiver.

(Proactively sending redundant symbols this way is called forward error
correction (FEC), and is useful to reduce or even eliminate round-trip delays
required for reliable object transmission.  It can be seen as a form of
“insurance,” with the amount of extra, redundant data being its “insurance
premium.”)

It is NOT okay for an encoding symbol to be corrupted—accidentally or
maliciously—during transit.  Feeding a Decoder with such corrupted symbols WILL
jeopardize recovery of the source object, so the receiver MUST detect and
discard corrupted encoding symbols, e.g. using checksums calculated and
transmitted by the sender along with the encoding symbol.

The Encoder and Decoder for the same source object should share the following
information: Total source object size (in octets), symbol size (chosen by the
sender), number of source blocks, number of sub-blocks (an internal detail), and
the symbol alignment factor (another internal detail).  The IETF Forward Error
Correction (FEC) Building Block specification (RFC 5052) and the RaptorQ
specification (RFC 6330) encapsulate these into two parameters: 64-bit Common
FEC Object Transmission Information, and 32-bit Scheme-Specific FEC Object
Transmission Information.  They are available from the sender's Encoder; the
receiver should pass them when creating a Decoder.
*/

package raptorq

// Encoder encodes one object into a series of symbols.
type Encoder interface {
	// Encode writes the encoding symbol identified by the given source block
	// number and encoding symbol ID, into the given buffer.
	//
	// On success, Encode returns the number of octets written into buf and nil
	// error.
	//
	// On error, Encode returns a non-nil error code.
	Encode(sbn uint8, esi uint32, buf []byte) (written uint, err error)

	// CommonOTI returns the Common FEC Object Transmission Information.
	CommonOTI() uint64

	// SchemeSpecificOTI returns the RaptorQ Scheme-Specific FEC Object
	// Transmission Information.
	SchemeSpecificOTI() uint32

	// TransferLength returns the source object size, in octets.  “F” in RFC
	// 6330.
	TransferLength() uint64

	// SymbolSize returns the symbol size, in octets.  “T” in RFC 6330.
	SymbolSize() uint16

	// NumSourceBlocks returns the number of source blocks.  “Z” in RFC 6330.
	NumSourceBlocks() uint8

	// NumSubBlocks returns the number of sub-blocks.  “N” in RFC 6330.
	NumSubBlocks() uint8

	// SymbolAlignmentParameter returns the symbol alignment, in octets.  “Al”
	// in RFC 6330.
	SymbolAlignmentParameter() uint8

	// MaxSubBlockSize returns the maximum size block that is decodable in
	// working memory, in octets.  “WS” in RFC 6330.
	MaxSubBlockSize() uint32

	// NumBlocks returns the number of blocks.
	NumBlocks() uint8

	// Free, on supported implementations, will free memory used for generating
	// encoding symbols for the given source block.  Once a source block has
	// been freed, calling Encode with its SBN may return an error.
	FreeSourceBlock(sbn uint8)

	// SourceBlockSize returns the size of the given source block, in bytes, or
	// 0 if sbn is out of range.
	SourceBlockSize(sbn uint8) uint32

	// NumSourceSymbols returns the number of encoding symbols covering the
	// original source data in the given source block, or 0 if sbn is out of
	// range.  The first encoding symbols up to this number is called source
	// symbols, and contains the source block data itself.
	NumSourceSymbols(sbn uint8) uint16

	// MinSymbols returns the minimum number of encoding symbols needed to
	// achieve the 99% probability of decoding the given source block,
	// or 0 if sbn is out of range.  “K” in RFC 6330.
	MinSymbols(sbn uint8) uint16

	// MaxSymbols returns the maximum number of encoding symbols the encoder
	// can generate for the given source block, or 0 if sbn is out of range.
	MaxSymbols(sbn uint8) uint32

	// Close closes the Encoder.  After an Encoder is closed, all methods but
	// Close() will panic if called.
	Close() error
}

// EncoderFactory is a factory of Encoder instances.
type EncoderFactory interface {
	/*
	NewEncoder creates and returns an Encoder that can encode the given source
	object into symbols.

	input is the source object to encode.

	symbolSize is the encoding symbol size, in octets.

	minSubSymbolSize is the minimum encoding symbol size allowed, in octets.
	(If you are not sure, or if you do not want internal interleaving of source
	symbols, set it equal to symbolSize.)

	maxSubBlockSize is the maximum size block that is decodable in working
	memory, in octets.  Iff this is lower than the source object size, the
    source object will be split into more than one source blocks.  The maximum
    allowed value is 56403 * symbolSize.

	alignment is an internal alignment parameter, in bytes.  Typically this is a
	power of 2, up to implementation-defined maximum.  Larger alignment will
	speed up calculation, at the expense of slightly higher transmission size
	overhead.  Both symbolSize and minSubSymbolSize must be a multiple of this.

	On success, NewEncoder returns an Encoder instance and nil error; on
    failure, it returns nil Encoder and an error code.
	*/
	New(input []byte, symbolSize uint16, minSubSymbolSize uint16,
        maxSubBlockSize uint32, alignment uint8) (enc Encoder, err error)
}

// TODO ek - below

// Decoder decodes encoding symbols and reconstructs one object into a series of
// symbols.
type Decoder interface {
}

type DecoderFactory interface {
}