%{
#include <stdint.h>
#include <cstdint>
%}

%include "stdint.i"

namespace std {
	using uint8_t = ::uint8_t;
	using int8_t = ::int8_t;
	using uint16_t = ::uint16_t;
	using int16_t = ::int16_t;
	using uint32_t = ::uint32_t;
	using int32_t = ::int32_t;
	using uint64_t = ::uint64_t;
	using int64_t = ::int64_t;
}

%typemap(gotype) uint8_t, uint8_t const, uint8_t const &, std::uint8_t, std::uint8_t const, std::uint8_t const & "uint8"
%typemap(gotype) int8_t, int8_t const, int8_t const &, std::int8_t, std::int8_t const, std::int8_t const & "int8"
%typemap(gotype) uint16_t, uint16_t const, uint16_t const &, std::uint16_t, std::uint16_t const, std::uint16_t const & "uint16"
%typemap(gotype) int16_t, int16_t const, int16_t const &, std::int16_t, std::int16_t const, std::int16_t const & "int16"
%typemap(gotype) uint32_t, uint32_t const, uint32_t const &, std::uint32_t, std::uint32_t const, std::uint32_t const & "uint32"
%typemap(gotype) int32_t, int32_t const, int32_t const &, std::int32_t, std::int32_t const, std::int32_t const & "int32"
%typemap(gotype) uint64_t, uint64_t const, uint64_t const &, std::uint64_t, std::uint64_t const, std::uint64_t const & "uint64"
%typemap(gotype) int64_t, int64_t const, int64_t const &, std::int64_t, std::int64_t const, std::int64_t const & "int64"
