package main

import (
	"io"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/ipld/go-ipld-prime/schema"
	gengo "github.com/ipld/go-ipld-prime/schema/gen/go"
)

type typedNodeGenerator interface {

	// -- the natively-typed apis -->
	//   (might be more readable to group these in another interface and have it
	//     return a `typedNodeGenerator` with the rest?  but structurally same.)

	EmitNativeType(io.Writer)
	EmitNativeAccessors(io.Writer) // depends on the kind -- field accessors for struct, typed iterators for map, etc.
	EmitNativeBuilder(io.Writer)   // typically emits some kind of struct that has a Build method.
	EmitNativeMaybe(io.Writer)     // a pointer-free 'maybe' mechanism is generated for all types.

	// -- the schema.TypedNode.Type method and vars -->

	EmitTypedNodeMethodType(io.Writer) // these emit dummies for now

	// -- all node methods -->
	//   (and note that the nodeBuilder for this one should be the "semantic" one,
	//     e.g. it *always* acts like a map for structs, even if the repr is different.)

	nodeGenerator

	// -- and the representation and its node and nodebuilder -->

	EmitTypedNodeMethodRepresentation(io.Writer)
}

type typedLinkNodeGenerator interface {
	// all methods in typedNodeGenerator
	typedNodeGenerator

	// as typed.LinkNode.ReferencedNodeBuilder generator
	EmitTypedLinkNodeMethodReferencedNodeBuilder(io.Writer)
}

type nodeGenerator interface {
	EmitNodeType(io.Writer)
	EmitNodeMethodReprKind(io.Writer)
	EmitNodeMethodLookupString(io.Writer)
	EmitNodeMethodLookup(io.Writer)
	EmitNodeMethodLookupIndex(io.Writer)
	EmitNodeMethodLookupSegment(io.Writer)
	EmitNodeMethodMapIterator(io.Writer)  // also iterator itself
	EmitNodeMethodListIterator(io.Writer) // also iterator itself
	EmitNodeMethodLength(io.Writer)
	EmitNodeMethodIsUndefined(io.Writer)
	EmitNodeMethodIsNull(io.Writer)
	EmitNodeMethodAsBool(io.Writer)
	EmitNodeMethodAsInt(io.Writer)
	EmitNodeMethodAsFloat(io.Writer)
	EmitNodeMethodAsString(io.Writer)
	EmitNodeMethodAsBytes(io.Writer)
	EmitNodeMethodAsLink(io.Writer)
}

func emitEntireType(ng nodeGenerator, w io.Writer) {
	if ng == nil {
		return
	}
	ng.EmitNodeType(w)
	ng.EmitNodeMethodReprKind(w)
	ng.EmitNodeMethodLookupString(w)
	ng.EmitNodeMethodLookup(w)
	ng.EmitNodeMethodLookupIndex(w)
	ng.EmitNodeMethodLookupSegment(w)
	ng.EmitNodeMethodMapIterator(w)
	ng.EmitNodeMethodListIterator(w)
	ng.EmitNodeMethodLength(w)
	ng.EmitNodeMethodIsUndefined(w)
	ng.EmitNodeMethodIsNull(w)
	ng.EmitNodeMethodAsBool(w)
	ng.EmitNodeMethodAsInt(w)
	ng.EmitNodeMethodAsFloat(w)
	ng.EmitNodeMethodAsString(w)
	ng.EmitNodeMethodAsBytes(w)
	ng.EmitNodeMethodAsLink(w)

	tg, ok := ng.(typedNodeGenerator)
	if ok {
		tg.EmitNativeType(w)
		tg.EmitNativeAccessors(w)
		tg.EmitNativeBuilder(w)
		tg.EmitNativeMaybe(w)
		tg.EmitTypedNodeMethodType(w)
		tg.EmitTypedNodeMethodRepresentation(w)
	}
	tlg, ok := ng.(typedLinkNodeGenerator)
	if ok {
		tlg.EmitTypedLinkNodeMethodReferencedNodeBuilder(w)
	}
}

func main() {
	openOrPanic := func(filename string) *os.File {
		y, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		return y
	}

	tString := schema.SpawnString("String")
	tInt := schema.SpawnInt("Int")
	tLink := schema.SpawnLink("Link")
	tBytes := schema.SpawnBytes("Bytes")

	tPBLink := schema.SpawnStruct("PBLink",
		[]schema.StructField{
			schema.SpawnStructField("Hash", tLink, true, false),
			schema.SpawnStructField("Name", tString, true, false),
			schema.SpawnStructField("Tsize", tInt, true, false),
		},
		schema.StructRepresentation_Map{},
	)

	tPBLinks := schema.SpawnList("PBLinks", tPBLink, false)

	tPBNode := schema.SpawnStruct("PBNode",
		[]schema.StructField{
			schema.SpawnStructField("Links", tPBLinks, false, false),
			schema.SpawnStructField("Data", tBytes, false, false),
		},
		schema.StructRepresentation_Map{},
	)

	tRaw := schema.SpawnBytes("RawNode")

	f := openOrPanic("common_gen.go")
	gengo.EmitMinima("dagpb", f)

	f = openOrPanic("pb_node_gen.go")
	gengo.EmitFileHeader("dagpb", f)
	tg := gengo.NewGeneratorForKindString(tString)
	emitEntireType(tg, f)
	emitEntireType(tg.GetRepresentationNodeGen(), f)
	tg = gengo.NewGeneratorForKindInt(tInt)
	emitEntireType(tg, f)
	emitEntireType(tg.GetRepresentationNodeGen(), f)
	tg = gengo.NewGeneratorForKindBytes(tBytes)
	emitEntireType(tg, f)
	emitEntireType(tg.GetRepresentationNodeGen(), f)
	tg = gengo.NewGeneratorForKindLink(tLink)
	emitEntireType(tg, f)
	emitEntireType(tg.GetRepresentationNodeGen(), f)
	tg = gengo.NewGeneratorForKindStruct(tPBLink)
	emitEntireType(tg, f)
	emitEntireType(tg.GetRepresentationNodeGen(), f)
	tg = gengo.NewGeneratorForKindList(tPBLinks)
	emitEntireType(tg, f)
	emitEntireType(tg.GetRepresentationNodeGen(), f)
	tg = gengo.NewGeneratorForKindStruct(tPBNode)
	emitEntireType(tg, f)
	emitEntireType(tg.GetRepresentationNodeGen(), f)
	if err := f.Close(); err != nil {
		panic(err)
	}
	read, err := ioutil.ReadFile("pb_node_gen.go")
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile("ipld\\.ErrInvalidKey\\{[^}]*\\}")
	newContents := re.ReplaceAll(read, []byte("err"))

	err = ioutil.WriteFile("pb_node_gen.go", []byte(newContents), 0)
	if err != nil {
		panic(err)
	}

	f = openOrPanic("raw_node_gen.go")
	gengo.EmitFileHeader("dagpb", f)
	tg = gengo.NewGeneratorForKindBytes(tRaw)
	emitEntireType(tg, f)
	emitEntireType(tg.GetRepresentationNodeGen(), f)
}
