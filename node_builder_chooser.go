package dagpb

import (
	ipld "github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/traversal"
)

// AddDagPBSupportToChooser takes an existing NodeBuilderChooser and subs in
// Protobuf and Raw node builders where neccesary
func AddDagPBSupportToChooser(existing traversal.LinkTargetNodePrototypeChooser) traversal.LinkTargetNodePrototypeChooser {
	return func(lnk ipld.Link, lnkCtx ipld.LinkContext) (ipld.NodePrototype, error) {
		c, ok := lnk.(cidlink.Link)
		if !ok {
			return existing(lnk, lnkCtx)
		}
		switch c.Cid.Prefix().Codec {
		case 0x70:
			return _PBNode__NodePrototype{}, nil
		case 0x55:
			return _RawNode__NodePrototype{}, nil
		default:
			return existing(lnk, lnkCtx)
		}
	}
}
