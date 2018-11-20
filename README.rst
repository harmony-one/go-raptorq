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

``go-raptorq`` contains two main interfaces, ``Encoder`` and ``Decoder``.

In order to send a binary object, termed **source object,** the sender creates
an ``Encoder`` for the source object.  The ``Encoder`` first splits the source
object into one or more contiguous blocks, termed **source blocks,** then for
each source block, the ``Encoder`` generates encoded and serially numbered
chunks of data, termed **encoding symbols.**  For each source block, the sender
uses the ``Encoder`` to generate as many encoding symbols as needed by the
receiver to recover the source block.  The sender chooses the size of the
encoding symbol, in octets, when creating an ``Encoder``.

On the other side, the receiver creates a ``Decoder`` for the same source
object, then keeps feeding the ``Decoder`` with the encoding symbols received
from the sender, until the ``Decoder`` is able to recover the source object from
the encoding symbols.  Once the source object has been recovered, the receiver
``Close()``-s and discards the ``Decoder``.

An encoding symbol is a unit data of transmission.  Each encoding symbol is
identified by a (``sbn``, ``esi``) pair, where:

* ``sbn``, or **source block number,** is the 8-bit, zero-based serial number of
  the source block from which the symbol was generated;
* ``esi``, or **encoding symbol identifier** is the 24-bit, zero-based serial
  identifier of the encoding symbol within the source block.

Upon receipt of an encoding symbol, the receiver first needs to identify the
``Decoder`` to use for the source object to which the received encoding symbol
belongs, then feed the ``Decoder`` with the encoding symbol, along with its
``sbn`` and ``esi``.

In order to recover each source block with high probability, the receiver needs
as many encoding symbols as needed to fulfill the original size (in octets) of
the source block.  For example, a 1.8MB source block with 1800-byte encoding
symbols, the receiver would need 1000 encoding symbols with 99% probability.
Each additional symbol adds roughly “two nines” to the decoding success
probability.  In the example above:

* 1000 encoding symbols would mean 99.99%,
* 1001 encoding symbols would mean 99.9999%,
* 1002 encoding symbols would mean 99.999999%, and so on.

It is okay for an encoding symbol to be completely lost (erased) during transit.
For each encoding symbol lost, the encoder simply needs to generate and send
another encoding symbol.  The sender need not know which encoding symbol was
lost; a brand new encoding symbol would be able to replace for any previously
sent-and-lost encoding symbol.  The sender may even anticipate losses and send
additional encoding symbols in advance without having to wait for negative
acknowledgements (NAKs) from the receiver.

(Proactively sending redundant symbols this way is called forward error
correction (FEC), and is useful to reduce or even eliminate round-trip delays
required for reliable object transmission.  It can be seen as a form of
“insurance,” with the amount of extra, redundant data being its “insurance
premium.”)

It is NOT okay for an encoding symbol to be corrupted—accidentally or
maliciously—during transit.  Feeding a ``Decoder`` with such corrupted symbols
WILL jeopardize recovery of the source object, so the receiver MUST detect and
discard corrupted encoding symbols, e.g. using checksums calculated and
transmitted by the sender along with the encoding symbol.

The ``Encoder`` and ``Decoder`` for the same source object should share the
following information: Total source object size (in octets), symbol size (chosen
by the sender), number of source blocks, number of sub-blocks (an internal
detail), and the symbol alignment factor (another internal detail).  The IETF
Forward Error Correction (FEC) Building Block specification (RFC 5052) and the
RaptorQ specification (RFC 6330) encapsulate these into two parameters: 64-bit
Common FEC Object Transmission Information, and 32-bit Scheme-Specific FEC
Object Transmission Information.  They are available from the sender's
``Encoder``; the receiver should pass them when creating a ``Decoder``.

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
