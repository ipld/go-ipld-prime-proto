package dagpb_test

import (
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	ipld "github.com/ipld/go-ipld-prime"
	dagpb "github.com/ipld/go-ipld-prime-proto"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/ipld/go-ipld-prime/traversal"
	mh "github.com/multiformats/go-multihash"
	. "github.com/warpfork/go-wish"
)

func TestNodeBuilderChooser(t *testing.T) {
	nb1 := basicnode.Style__Any{}
	nb2 := basicnode.Style__String{}
	var nb1Chooser traversal.LinkTargetNodeStyleChooser = dagpb.AddDagPBSupportToChooser(func(ipld.Link, ipld.LinkContext) (ipld.NodeStyle, error) {
		return nb1, nil
	})
	var nb2Chooser traversal.LinkTargetNodeStyleChooser = dagpb.AddDagPBSupportToChooser(func(ipld.Link, ipld.LinkContext) (ipld.NodeStyle, error) {
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

	ns, err := nb1Chooser(protoLink, ipld.LinkContext{})
	Wish(t, err, ShouldEqual, nil)
	Wish(t, ns, ShouldEqual, dagpb.Style.Protobuf)
	ns, err = nb1Chooser(rawLink, ipld.LinkContext{})
	Wish(t, err, ShouldEqual, nil)
	Wish(t, ns, ShouldEqual, dagpb.Style.Raw)
	ns, err = nb1Chooser(cborLink, ipld.LinkContext{})
	Wish(t, err, ShouldEqual, nil)
	Wish(t, ns, ShouldEqual, nb1)
	ns, err = nb2Chooser(protoLink, ipld.LinkContext{})
	Wish(t, err, ShouldEqual, nil)
	Wish(t, ns, ShouldEqual, dagpb.Style.Protobuf)
	ns, err = nb2Chooser(rawLink, ipld.LinkContext{})
	Wish(t, err, ShouldEqual, nil)
	Wish(t, ns, ShouldEqual, dagpb.Style.Raw)
	ns, err = nb2Chooser(cborLink, ipld.LinkContext{})
	Wish(t, err, ShouldEqual, nil)
	Wish(t, ns, ShouldEqual, nb2)

}
