package pkg

type Integer = uint32

const (
	// Symbol alignment parameter, in octets (4 for uint32; see Integer above)
	Al = 4

	// K'_max: the maximum number of source symbols per source block
	MaxKp = 56403 // section 4.3, section 5.1.2
)

type Params struct {
	// F: the transfer length of the object, in octets
	F uint64

	// WS: the maximum size block that is decodable in working memory, in octets
	WS uint32

	// P': the maximum payload size in octets, which is assumed to be a multiple
	// of Al
	Pp uint32

	// SS: a parameter where the desired lower bound on the sub-symbol size is
	// SS*Al
	SS uint32
}

// A row in Table 2 in section 5.6
type Table2Row struct {
	Kp uint16
	JKp uint16
	SKp uint16
	HKp uint16
	LKp uint16
}

var Table2 Table2Row[...] = {
}

// T() returns the symbol size in octets, which MUST be a multiple of Al.
func (params Params) T() uint32 {
	return params.Pp
}

// Kt() returns number of symbols in the transfer object.
func (params Params) Kt() uint64 {
	F := params.F
	T := uint64(params.T())
	return (F + T - 1) / T
}

// MaxN() returns the maximum number of sub-symbols per symbol.
func (params Params) MaxN() uint32 {
	return params.T() / (params.SS * Al)
}

// KL(n) returns the largest possible Kp value that can fit into working memory.
func (params Params) KL(n uint32) uint16 {
	thresh := params.WS / (Al * ((params.T() + Al * n - 1) / (Al * n)))
}

Encoder = struct {

}
