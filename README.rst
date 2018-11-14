What Is ``go-raptorq``?
=======================
``go-raptorq`` implements the `RaptorQ`_ error correction code as defined in
`RFC 6330`_.  It includes a RFC-compliant encoder and decoder.

What Is RaptorQ?
================

RaptorQ is a rateless erasure code (also known as `fountain code`), and provides
two functions:

* Redundantly encoding a message into practically infinite number (~2**24) of
  *symbols*;
* Reliably decoding the original message from *any* subset of encoded symbols
  *with high probability*, provided that the cumulative size of the symbols
  received is equal to or slightly greater than the original message size.

RaptorQ is useful for a variety of purposes, including but not limited to:

* Transmitting a message reliably over a lossy and/or adversarial network path,
  without employing acknowledgement (feedback) mechanism or suffering from
  round-trip delays incurred thereby.
* Reliable object storage, where redundancy/fault-tolerence level—such as number
  of parity disks in RAID arrays—can be scaled up or down on demand without
  having to re-code contents on the existing disks.

We at `Harmony`_ are developing and using ``go-raptorq`` in order to implement a
near-optimal, adversary-resilient, and stable-latency message broadcasting
mechanism for use in our highly scalable and performant blockchain network.

Using ``go-raptorq``
====================

To use ``go-raptorq`` in your Go application::

  $ CGO_CXXFLAGS='-std=c++11' go get simple-rules/go-raptorq

``go-raptorq`` has two main structs, Encoder and Decoder.  One Encoder/Decoder
object is used to encode/decode one message.  See the godoc **(TODO)** for
details.

Licensing
=========

Copyright © 2018, Simple Rules Company.  All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

The views and conclusions contained in the software and documentation are those
of the authors and should not be interpreted as representing official policies,
either expressed or implied, of the go-raptorq project.

In addition to the terms and conditions of the license above, the licensee shall
comply with the terms and conditions associated with all `IETF IPR disclosure
associated with RFC 6330`_.  In no event shall the copyright owner or
contributors be liable for damages arising in any way out of failure to comply
with the terms and conditions of the above mentioned IPR disclosure.

.. _RaptorQ: https://www.qualcomm.com/media/documents/files/raptorq-technical-overview.pdf
.. _RFC 6330: https://tools.ietf.org/html/rfc6330
.. _IETF IPR disclosure associated with RFC 6330: https://datatracker.ietf.org/ipr/search/?rfc=6330&submit=rfc
.. _IETF IPR Disclosure ID #2554: https://datatracker.ietf.org/ipr/2554/
.. _fountain code: https://en.wikipedia.org/wiki/Fountain_code
.. _Harmony: https://harmony.one/
