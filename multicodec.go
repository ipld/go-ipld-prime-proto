package dagpb

import (
	"errors"
	"io"

	ipld "github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
)

var (
	ErrNoAutomaticDecoding                           = errors.New("No automatic decoding for this type, node builder must provide fast path")
	ErrNoAutomaticEncoding                           = errors.New("No automatic encoding for this type, node must provide fast path")
	_                      cidlink.MulticodecDecoder = DagPBDecoder
	_                      cidlink.MulticodecEncoder = DagPBEncoder
	_                      cidlink.MulticodecDecoder = RawDecoder
	_                      cidlink.MulticodecEncoder = RawEncoder
)

func init() {
	cidlink.RegisterMulticodecDecoder(0x70, DagPBDecoder)
	cidlink.RegisterMulticodecEncoder(0x70, DagPBEncoder)
	cidlink.RegisterMulticodecDecoder(0x55, RawDecoder)
	cidlink.RegisterMulticodecEncoder(0x55, RawEncoder)
}

func DagPBDecoder(nb ipld.NodeBuilder, r io.Reader) (ipld.Node, error) {
	// Probe for a builtin fast path.  Shortcut to that if possible.
	//  (ipldcbor.NodeBuilder supports this, for example.)
	type detectFastPath interface {
		DecodeDagProto(io.Reader) (ipld.Node, error)
	}
	if nb2, ok := nb.(detectFastPath); ok {
		return nb2.DecodeDagProto(r)
	}
	// Okay, generic builder path.
	return nil, ErrNoAutomaticDecoding
}

func DagPBEncoder(n ipld.Node, w io.Writer) error {
	// Probe for a builtin fast path.  Shortcut to that if possible.
	//  (ipldcbor.Node supports this, for example.)
	type detectFastPath interface {
		EncodeDagProto(io.Writer) error
	}
	if n2, ok := n.(detectFastPath); ok {
		return n2.EncodeDagProto(w)
	}
	// Okay, generic inspection path.
	return ErrNoAutomaticEncoding
}

func RawDecoder(nb ipld.NodeBuilder, r io.Reader) (ipld.Node, error) {
	// Probe for a builtin fast path.  Shortcut to that if possible.
	//  (ipldcbor.NodeBuilder supports this, for example.)
	type detectFastPath interface {
		DecodeDagRaw(io.Reader) (ipld.Node, error)
	}
	if nb2, ok := nb.(detectFastPath); ok {
		return nb2.DecodeDagRaw(r)
	}
	// Okay, generic builder path.
	return nil, ErrNoAutomaticDecoding
}

func RawEncoder(n ipld.Node, w io.Writer) error {
	// Probe for a builtin fast path.  Shortcut to that if possible.
	//  (ipldcbor.Node supports this, for example.)
	type detectFastPath interface {
		EncodeDagRaw(io.Writer) error
	}
	if n2, ok := n.(detectFastPath); ok {
		return n2.EncodeDagRaw(w)
	}
	// Okay, generic inspection path.
	return ErrNoAutomaticEncoding
}
