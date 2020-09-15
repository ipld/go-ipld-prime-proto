package main

import (
	"io/ioutil"
	"os"
	"regexp"

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

	adjCfg := &gengo.AdjunctCfg{}

	pkgName := "dagpb"

	f := openOrPanic("common_gen.go")
	gengo.EmitInternalEnums(pkgName, f)

	f = openOrPanic("pb_node_gen.go")
	gengo.EmitFileHeader(pkgName, f)
	gengo.EmitEntireType(gengo.NewStringReprStringGenerator(pkgName, tString, adjCfg), f)
	gengo.EmitEntireType(gengo.NewIntReprIntGenerator(pkgName, tInt, adjCfg), f)
	gengo.EmitEntireType(gengo.NewBytesReprBytesGenerator(pkgName, tBytes, adjCfg), f)
	gengo.EmitEntireType(gengo.NewLinkReprLinkGenerator(pkgName, tLink, adjCfg), f)
	gengo.EmitEntireType(gengo.NewStructReprMapGenerator(pkgName, tPBLink, adjCfg), f)
	gengo.EmitEntireType(gengo.NewListReprListGenerator(pkgName, tPBLinks, adjCfg), f)
	gengo.EmitEntireType(gengo.NewStructReprMapGenerator(pkgName, tPBNode, adjCfg), f)

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
	gengo.EmitEntireType(gengo.NewBytesReprBytesGenerator(pkgName, tRaw, adjCfg), f)
}
