%{
#include <stddef.h>
%}

/*
 * SLICE_TYPEMAP enables passing of a slice as an argument into a function which
 * accepts the following buffer interfaces.  The following C/C++ argument
 * patterns are recognized:
 *
 * TYPE *SLICE, size_t SLICELEN
 * TYPE *SLICEBEGIN, char *SLICEEND
 * TYPE *&SLICEBEGIN, char *&SLICEEND
 *
 * The first form is the typical (array pointer, number of elements) convention
 * used mostly in C programs.  The second form is the (start pointer, end
 * pointer) form compatible with the begin/end iterator convention in C++.  The
 * third form is a slight variation of the second form, where the function moves
 * the begin iterator passed by reference in order to let caller know how many
 * elements have been consumed or filled.
 *
 * []byte slices can be passed for char * or unsigned char *; []int8 slices can
 * be passed for signed char *.  Users can also introduce other C-Go type
 * mappings by using SLICE_TYPEMAP(C type, Go type).
 */
%define SLICE_TYPEMAP(TYPE, GOTYPE)
%typemap(gotype) (TYPE *SLICE, size_t SLICELEN), (TYPE *SLICEBEGIN, TYPE *SLICEEND), (TYPE *&SLICEBEGIN, TYPE *SLICEEND) %{[]GOTYPE%}

%typemap(in) (TYPE *SLICE, size_t SLICELEN)
{
  $1 = (TYPE *)$input.array;
  $2 = $input.len;
}
%typemap(out) (TYPE *SLICE, size_t SLICELEN) ""
%typemap(goout) (TYPE *SLICE, size_t SLICELEN) ""
%typemap(freearg) (TYPE *SLICE, size_t SLICELEN) ""
%typemap(argout) (TYPE *SLICE, size_t SLICELEN) ""

%typemap(in) (TYPE *SLICEBEGIN, TYPE *SLICEEND)
{
  $1 = (TYPE *)$input.array;
  $2 = $1 + $input.len;
}
%typemap(out) (TYPE *SLICEBEGIN, TYPE *SLICEEND) ""
%typemap(goout) (TYPE *SLICEBEGIN, TYPE *SLICEEND) ""
%typemap(freearg) (TYPE *SLICEBEGIN, TYPE *SLICEEND) ""
%typemap(argout) (TYPE *SLICEBEGIN, TYPE *SLICEEND) ""

%typemap(in) (TYPE *&SLICEBEGIN, TYPE *SLICEEND) (TYPE *temp)
{
  temp = (TYPE *)$input.array;
  $1 = &temp;
  $2 = temp + $input.len;
}
%typemap(out) (TYPE *&SLICEBEGIN, TYPE *SLICEEND) ""
%typemap(goout) (TYPE *&SLICEBEGIN, TYPE *SLICEEND) ""
%typemap(freearg) (TYPE *&SLICEBEGIN, TYPE *SLICEEND) ""
%typemap(argout) (TYPE *&SLICEBEGIN, TYPE *SLICEEND) ""

%enddef

SLICE_TYPEMAP(signed char, int8);
SLICE_TYPEMAP(unsigned char, byte);
SLICE_TYPEMAP(char, byte);
