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
	var nb1Chooser traversal.NodeBuilderChooser = dagpb.AddDagPBSupportToChooser(func(ipld.Link, ipld.LinkContext) (ipld.NodeBuilder, error) {
		return nb1, nil
	})
	var nb2Chooser traversal.NodeBuilderChooser = dagpb.AddDagPBSupportToChooser(func(ipld.Link, ipld.LinkContext) (ipld.NodeBuilder, error) {
		return nb2, nil
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

	nb, err := nb1Chooser(protoLink, ipld.LinkContext{})
	Wish(t, err, ShouldEqual, nil)
	Wish(t, nb, ShouldEqual, dagpb.PBNode__NodeBuilder())
	nb, err = nb1Chooser(rawLink, ipld.LinkContext{})
	Wish(t, err, ShouldEqual, nil)
	Wish(t, nb, ShouldEqual, dagpb.RawNode__NodeBuilder())
	nb, err = nb1Chooser(cborLink, ipld.LinkContext{})
	Wish(t, err, ShouldEqual, nil)
	Wish(t, nb, ShouldEqual, nb1)
	nb, err = nb2Chooser(protoLink, ipld.LinkContext{})
	Wish(t, err, ShouldEqual, nil)
	Wish(t, nb, ShouldEqual, dagpb.PBNode__NodeBuilder())
	nb, err = nb2Chooser(rawLink, ipld.LinkContext{})
	Wish(t, err, ShouldEqual, nil)
	Wish(t, nb, ShouldEqual, dagpb.RawNode__NodeBuilder())
	nb, err = nb2Chooser(cborLink, ipld.LinkContext{})
	Wish(t, err, ShouldEqual, nil)
	Wish(t, nb, ShouldEqual, nb2)

}
