package main

import (
	"os"

	"github.com/ipld/go-ipld-prime/schema"
	gengo "github.com/ipld/go-ipld-prime/schema/gen/go"
)

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
	gengo.EmitEntireType(gengo.NewGeneratorForKindString(tString), f)
	gengo.EmitEntireType(gengo.NewGeneratorForKindInt(tInt), f)
	gengo.EmitEntireType(gengo.NewGeneratorForKindBytes(tBytes), f)
	gengo.EmitEntireType(gengo.NewGeneratorForKindLink(tLink), f)
	gengo.EmitEntireType(gengo.NewGeneratorForKindStruct(tPBLink), f)
	gengo.EmitEntireType(gengo.NewGeneratorForKindList(tPBLinks), f)
	gengo.EmitEntireType(gengo.NewGeneratorForKindStruct(tPBNode), f)

	f = openOrPanic("raw_node_gen.go")
	gengo.EmitFileHeader("dagpb", f)
	gengo.EmitEntireType(gengo.NewGeneratorForKindBytes(tRaw), f)
}
