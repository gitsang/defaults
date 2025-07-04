package main

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/gitsang/defaults"
)

type Gender string

type Sample struct {
	Name    string `default:"John Smith"`
	Age     int    `default:"27"`
	Gender  Gender `default:"m"`
	Working bool   `default:"true"`

	SliceInt    []int    `default:"[1, 2, 3]"`
	SlicePtr    []*int   `default:"[1, 2, 3]"`
	SliceString []string `default:"[\"a\", \"b\"]"`

	MapNull            map[string]int          `default:"{}"`
	Map                map[string]int          `default:"{\"key1\": 123}"`
	MapOfStruct        map[string]OtherStruct  `default:"{\"Key2\": {\"Foo\":123}}"`
	MapOfPtrStruct     map[string]*OtherStruct `default:"{\"Key3\": {\"Foo\":123}}"`
	MapOfStructWithTag map[string]OtherStruct  `default:"{\"Key4\": {\"Foo\":123}}"`

	MapOfStructWithOutTag map[string]OtherStruct
	MapOfSliceWithOutTag  map[string][]OtherStruct

	Struct    OtherStruct  `default:"{\"Foo\": 123}"`
	StructPtr *OtherStruct `default:"{\"Foo\": 123}"`

	NoTag    OtherStruct // Recurses into a nested struct by default
	NoOption OtherStruct `default:"-"` // no option
}

type OtherStruct struct {
	Hello  string `default:"world"` // Tags in a nested struct also work
	Foo    int    `default:"-"`
	Random int    `default:"-"`
}

// SetDefaults implements defaults.Setter interface
func (s *OtherStruct) SetDefaults() {
	if defaults.CanUpdate(s.Random) { // Check if it's a zero value (recommended)
		s.Random = rand.Int() // Set a dynamic value
	}
}

func main() {
	obj := &Sample{
		MapOfStructWithOutTag: map[string]OtherStruct{
			"hello": {},
		},
		MapOfSliceWithOutTag: map[string][]OtherStruct{
			"hello": {
				{},
			},
		},
	}
	if err := defaults.Set(obj); err != nil {
		panic(err)
	}

	out, err := json.MarshalIndent(obj, "", "	")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}
