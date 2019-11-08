package dagpb_test

import (
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	ipld "github.com/ipld/go-ipld-prime"
	dagpb "github.com/ipld/go-ipld-prime-proto"
	ipldfree "github.com/ipld/go-ipld-prime/impl/free"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/traversal"
	mh "github.com/multiformats/go-multihash"
	. "github.com/warpfork/go-wish"
)

func TestNodeBuilderChooser(t *testing.T) {
	nb1 := ipldfree.NodeBuilder()
	nb2 := dagpb.String__NodeBuilder()
	var nb1Chooser traversal.NodeBuilderChooser = dagpb.AddDagPBSupportToChooser(func(ipld.Link, ipld.LinkContext) ipld.NodeBuilder {
		return nb1
	})
	var nb2Chooser traversal.NodeBuilderChooser = dagpb.AddDagPBSupportToChooser(func(ipld.Link, ipld.LinkContext) ipld.NodeBuilder {
		return nb2
	})
	bytes := randomBytes(256)
	protoPrefix := merkledag.V1CidPrefix()
	protoCid, err := protoPrefix.Sum(bytes)
	Wish(t, err, ShouldEqual, nil)
	rawPrefix := cid.Prefix{
		Codec:    cid.Raw,
		MhLength: -1,
		MhType:   mh.SHA2_256,
		Version:  1,
	}
	rawCid, err := rawPrefix.Sum(bytes)
	Wish(t, err, ShouldEqual, nil)
	cborPrefix := cid.Prefix{
		Codec:    cid.DagCBOR,
		MhLength: -1,
		MhType:   mh.SHA2_256,
		Version:  1,
	}
	cborCid, err := cborPrefix.Sum(bytes)
	Wish(t, err, ShouldEqual, nil)

	protoLink := cidlink.Link{Cid: protoCid}
	rawLink := cidlink.Link{Cid: rawCid}
	cborLink := cidlink.Link{Cid: cborCid}

	Wish(t, nb1Chooser(protoLink, ipld.LinkContext{}), ShouldEqual, dagpb.PBNode__NodeBuilder())
	Wish(t, nb1Chooser(rawLink, ipld.LinkContext{}), ShouldEqual, dagpb.RawNode__NodeBuilder())
	Wish(t, nb1Chooser(cborLink, ipld.LinkContext{}), ShouldEqual, nb1)
	Wish(t, nb2Chooser(protoLink, ipld.LinkContext{}), ShouldEqual, dagpb.PBNode__NodeBuilder())
	Wish(t, nb2Chooser(rawLink, ipld.LinkContext{}), ShouldEqual, dagpb.RawNode__NodeBuilder())
	Wish(t, nb2Chooser(cborLink, ipld.LinkContext{}), ShouldEqual, nb2)

}
