# defaults

[![CircleCI](https://circleci.com/gh/gitsang/defaults/tree/master.svg?style=svg)](https://circleci.com/gh/gitsang/defaults/tree/master)
[![codecov](https://codecov.io/gh/gitsang/defaults/branch/master/graph/badge.svg)](https://codecov.io/gh/gitsang/defaults)
[![GitHub release](https://img.shields.io/github/release/gitsang/defaults.svg)](https://github.com/gitsang/defaults/releases)
[![License](https://img.shields.io/github/license/gitsang/defaults.svg)](./LICENSE)

Initialize structs with default values

- Supports almost all kind of types
  - Scalar types
    - `int/8/16/32/64`, `uint/8/16/32/64`, `float32/64`
    - `uintptr`, `bool`, `string`
  - Complex types
    - `map`, `slice`, `struct`
  - Nested types
    - `map[K1]map[K2]Struct`, `[]map[K1]Struct[]`
  - Aliased types
    - `time.Duration`
    - e.g., `type Enum string`
  - Pointer types
    - e.g., `*SampleStruct`, `*int`
- Recursively initializes fields in a struct
- Dynamically sets default values by [`defaults.Setter`](./setter.go) interface
- Preserves non-initial values from being reset with a default value
- Set default values for all addressable structs

## 1. Usage

```go
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
```

output

```json
{
  "Name": "John Smith",
  "Age": 27,
  "Gender": "m",
  "Working": true,
  "SliceInt": [1, 2, 3],
  "SlicePtr": [1, 2, 3],
  "SliceString": ["a", "b"],
  "MapNull": {},
  "Map": {
    "key1": 123
  },
  "MapOfStruct": {
    "Key2": {
      "Hello": "world",
      "Foo": 123,
      "Random": 6012378114984103869
    }
  },
  "MapOfPtrStruct": {
    "Key3": {
      "Hello": "world",
      "Foo": 123,
      "Random": 2848306665585403127
    }
  },
  "MapOfStructWithTag": {
    "Key4": {
      "Hello": "world",
      "Foo": 123,
      "Random": 7421038243474730390
    }
  },
  "MapOfStructWithOutTag": {
    "hello": {
      "Hello": "world",
      "Foo": 0,
      "Random": 6214942057906930174
    }
  },
  "MapOfSliceWithOutTag": {
    "hello": [
      {
        "Hello": "world",
        "Foo": 0,
        "Random": 5317591004601161060
      }
    ]
  },
  "Struct": {
    "Hello": "world",
    "Foo": 123,
    "Random": 3365636741323893358
  },
  "StructPtr": {
    "Hello": "world",
    "Foo": 123,
    "Random": 4271644362335896782
  },
  "NoTag": {
    "Hello": "world",
    "Foo": 0,
    "Random": 8353712922077708401
  },
  "NoOption": {
    "Hello": "",
    "Foo": 0,
    "Random": 0
  }
}
```
