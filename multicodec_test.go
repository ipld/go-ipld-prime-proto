package dagpb_test

import (
	"bytes"
	"testing"

	blocks "github.com/ipfs/go-block-format"
	ipld "github.com/ipld/go-ipld-prime"
	dagpb "github.com/ipld/go-ipld-prime-proto"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/jbenet/go-random"
	. "github.com/warpfork/go-wish"
)

var seedSeq int64

func randomBytes(n int64) []byte {
	data := new(bytes.Buffer)
	random.WritePseudoRandomBytes(n, data, seedSeq)
	seedSeq++
	return data.Bytes()
}

func makeRawNode(randBytes []byte) (ipld.Node, error) {
	raw_nb := dagpb.RawNode__NodeBuilder()
	return raw_nb.CreateBytes(randBytes)
}

func makeProtoNode(linkedNodes map[string]ipld.Node) (ipld.Node, error) {
	dagpb_nb := dagpb.PBNode__NodeBuilder()
	dagpb_mb, err := dagpb_nb.CreateMap()
	if err != nil {
		return nil, err
	}
	linksKey, err := dagpb_mb.BuilderForKeys().CreateString("Links")
	if err != nil {
		return nil, err
	}
	daglinks_nb := dagpb_mb.BuilderForValue("Links")
	daglinks_lb, err := daglinks_nb.CreateList()
	if err != nil {
		return nil, err
	}
	i := 0
	for name, node := range linkedNodes {
		daglink_nb := daglinks_lb.BuilderForValue(i)
		daglink_mb, err := daglink_nb.CreateMap()
		if err != nil {
			return nil, err
		}
		hashKey, err := daglink_mb.BuilderForKeys().CreateString("Hash")
		if err != nil {
			return nil, err
		}
		nameKey, err := daglink_mb.BuilderForKeys().CreateString("Name")
		if err != nil {
			return nil, err
		}
		tsizeKey, err := daglink_mb.BuilderForKeys().CreateString("Tsize")
		if err != nil {
			return nil, err
		}
		nodeBytes, err := node.AsBytes()
		if err != nil {
			return nil, err
		}
		blk := blocks.NewBlock(nodeBytes)
		hashNode, err := daglink_mb.BuilderForValue("Hash").CreateLink(cidlink.Link{Cid: blk.Cid()})
		if err != nil {
			return nil, err
		}
		err = daglink_mb.Insert(hashKey, hashNode)
		if err != nil {
			return nil, err
		}
		tsizeNode, err := daglink_mb.BuilderForValue("Tsize").CreateInt(len(nodeBytes))
		if err != nil {
			return nil, err
		}
		err = daglink_mb.Insert(tsizeKey, tsizeNode)
		if err != nil {
			return nil, err
		}
		nameNode, err := daglink_mb.BuilderForValue("Name").CreateString(name)
		if err != nil {
			return nil, err
		}
		err = daglink_mb.Insert(nameKey, nameNode)
		if err != nil {
			return nil, err
		}
		linkNode, err := daglink_mb.Build()
		if err != nil {
			return nil, err
		}
		err = daglinks_lb.Append(linkNode)
		if err != nil {
			return nil, err
		}
		i++
	}
	linksNode, err := daglinks_lb.Build()
	if err != nil {
		return nil, err
	}
	err = dagpb_mb.Insert(linksKey, linksNode)
	if err != nil {
		return nil, err
	}
	dataKey, err := dagpb_mb.BuilderForKeys().CreateString("Data")
	if err != nil {
		return nil, err
	}
	randBytes := randomBytes(1000)
	dataNode, err := dagpb_mb.BuilderForValue("Data").CreateBytes(randBytes)
	if err != nil {
		return nil, err
	}
	err = dagpb_mb.Insert(dataKey, dataNode)
	if err != nil {
		return nil, err
	}
	return dagpb_mb.Build()
}

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
