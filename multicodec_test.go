package dagpb_test

import (
	"bytes"
	"testing"

	ipld "github.com/ipld/go-ipld-prime"
	dagpb "github.com/ipld/go-ipld-prime-proto"
	. "github.com/warpfork/go-wish"
)

func TestRoundTripRaw(t *testing.T) {
	randBytes := randomBytes(256)
	rawNode, err := makeRawNode(randBytes)
	Wish(t, err, ShouldEqual, nil)
	t.Run("encoding", func(t *testing.T) {
		var buf bytes.Buffer
		err := dagpb.RawEncoder(rawNode, &buf)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, buf.Bytes(), ShouldEqual, randBytes)
	})
	t.Run("decoding", func(t *testing.T) {
		buf := bytes.NewBuffer(randBytes)
		rawNode2, err := dagpb.RawDecoder(dagpb.RawNode__NodeBuilder(), buf)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, rawNode2, ShouldEqual, rawNode)
	})
}

func TestRoundTripProtbuf(t *testing.T) {
	randBytes1 := randomBytes(256)
	rawNode1, err := makeRawNode(randBytes1)
	Wish(t, err, ShouldEqual, nil)
	randBytes2 := randomBytes(256)
	rawNode2, err := makeRawNode(randBytes2)
	Wish(t, err, ShouldEqual, nil)
	pbNode, err := makeProtoNode(map[string]ipld.Node{
		"applesuace": rawNode1,
		"oranges":    rawNode2,
	})
	Wish(t, err, ShouldEqual, nil)
	t.Run("encode/decode equivalency", func(t *testing.T) {
		var buf bytes.Buffer
		err := dagpb.DagPBEncoder(pbNode, &buf)
		Wish(t, err, ShouldEqual, nil)
		pbNode2, err := dagpb.DagPBDecoder(dagpb.PBNode__NodeBuilder(), &buf)
		Wish(t, err, ShouldEqual, nil)
		Wish(t, pbNode2, ShouldEqual, pbNode)
	})
}
