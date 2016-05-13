// This is a very basic example of a program that implements rdb.decoder and
// outputs a human readable diffable dump of the rdb file.
// copy from github.com/cupcake/rdb/examples/diff.go

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cupcake/rdb"
	"github.com/cupcake/rdb/nopdecoder"
	"os"
	"strings"
)

var version = "0.3 20160513"
var jsonType = flag.Bool("json", false, "print as json")

const (
	kTypeString    string = "string"
	kTypeHset             = "hset"
	kTypeSet              = "set"
	kTypeList             = "list"
	kTypeSortedSet        = "sortedset"
)

var allKeyTypes = map[string]int{
	kTypeString:    1,
	kTypeHset:      1,
	kTypeSet:       1,
	kTypeList:      1,
	kTypeSortedSet: 1,
}

var kTypePrint = map[string]bool{}

var printValue = false

type decoder struct {
	db int
	i  int
	nopdecoder.NopDecoder
}

func (p *decoder) StartDatabase(n int) {
	p.db = n
}

type jsonStruct map[string]interface{}

func (js jsonStruct) String() string {
	bs, _ := json.Marshal(js)
	return string(bs)
}

func (p *decoder) Set(key, value []byte, expiry int64) {
	if kTypePrint[kTypeString] {
		if printValue {
			if *jsonType {
				fmt.Println(jsonStruct{"type": kTypeString, "key_b": key, "value_b": value, "expiry": expiry, "value_len": len(value)})
			} else {
				fmt.Printf("%s\t%q\t%d\tvalue:\t%q\texpiry:\t%d\n", kTypeString, key, len(value), value, expiry)
			}
		} else {
			if *jsonType {
				fmt.Println(jsonStruct{"type": kTypeString, "key_b": key, "expiry": expiry, "value_len": len(value)})
			} else {
				fmt.Printf("%s\t%q\t%d\n", kTypeString, key, len(value))
			}
		}
	}
}

func (p *decoder) Hset(key, field, value []byte) {
	if kTypePrint[kTypeHset] {
		if printValue {
			if *jsonType {
				fmt.Println(jsonStruct{"type": kTypeHset, "key_b": key, "value_b": value, "field": field, "value_len": len(value)})
			} else {
				fmt.Printf("%s\t%q\t%d\tfield:\t%q\tvalue:\t%q\t\n", kTypeHset, key, len(field)+len(value), field, value)
			}
		} else {
			if *jsonType {
				fmt.Println(jsonStruct{"type": kTypeHset, "key_b": key, "field_b": field, "value_len": len(value)})
			} else {
				fmt.Printf("%s\t%q\t%d\n", kTypeHset, key, len(field)+len(value))
			}
		}
	}
}

func (p *decoder) Sadd(key, member []byte) {
	if kTypePrint[kTypeSet] {
		if printValue {
			if *jsonType {
				fmt.Println(jsonStruct{"type": kTypeSet, "key_b": key, "member_b": member, "value_len": len(member)})
			} else {
				fmt.Printf("%s\t%q\t%d\tmember:\t%q\n", kTypeSet, key, len(member), member)
			}
		} else {
			if *jsonType {
				fmt.Println(jsonStruct{"type": kTypeSet, "key_b": key, "value_len": len(member)})
			} else {
				fmt.Printf("%s\t%q\t%d\n", kTypeSet, key, len(member))
			}
		}
	}
}

func (p *decoder) StartList(key []byte, length, expiry int64) {
	p.i = 0
}

func (p *decoder) Rpush(key, value []byte) {
	if kTypePrint[kTypeList] {
		if printValue {
			if *jsonType {
				fmt.Println(jsonStruct{"type": kTypeList, "key_b": key, "value_b": value, "value_len": len(value)})
			} else {
				fmt.Printf("%s\t%q\t%d\tvalue:\t%q\n", kTypeList, key, len(value), value)
			}
		} else {
			if *jsonType {
				fmt.Println(jsonStruct{"type": kTypeList, "key_b": key, "value_len": len(value)})
			} else {
				fmt.Printf("%s\t%q\t%d\n", kTypeList, key, len(value))
			}
		}
	}
	p.i++
}

func (p *decoder) StartZSet(key []byte, cardinality, expiry int64) {
	p.i = 0
}

func (p *decoder) Zadd(key []byte, score float64, member []byte) {
	if kTypePrint[kTypeSortedSet] {
		if printValue {
			if *jsonType {
				fmt.Println(jsonStruct{"type": kTypeSortedSet, "key_b": key, "member_b": member, "score": score, "value_len": len(member)})
			} else {
				fmt.Printf("%s\t%s\t%d\tscore:\t%f\tmember:\t%q\n", kTypeSortedSet, string(key), len(member), score, member)
			}
		} else {
			if *jsonType {
				fmt.Println(jsonStruct{"type": kTypeSortedSet, "key_b": key, "score": score, "value_len": len(member)})
			} else {
				fmt.Printf("%s\t%s\t%d\n", kTypeSortedSet, string(key), len(member))
			}
		}
	}
	p.i++
}

func maybeFatal(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err)
		os.Exit(1)
	}
}

var printTypes string

var typesAll []string

func init() {
	for kt := range allKeyTypes {
		kTypePrint[kt] = false
		typesAll = append(typesAll, kt)
	}
	flag.StringVar(&printTypes, "types", "", "print these types:["+strings.Join(typesAll, ",")+"],default all")
	flag.BoolVar(&printValue, "val", false, "show values")
	df := flag.Usage
	flag.Usage = func() {
		df()
		name := os.Args[0]
		fmt.Println("\n redis rdb viewer")
		fmt.Println(" " + name + " -types string part1.rdb")
		fmt.Println(" " + name + " -types string,set -val part1.rdb")
		fmt.Println("")
		fmt.Println(` first three parts : [string "abc" 12] -->[type key value_len]`)
		fmt.Println("\n site: https://github.com/hidu/rdb-viewer")
		fmt.Println(" version:", version)
	}
}

func main() {
	flag.Parse()

	fileName := flag.Arg(0)

	if fileName == "" {
		fmt.Println("rdb file name required")
		os.Exit(2)
	}
	parseArgTypes()

	f, err := os.Open(fileName)
	maybeFatal(err)
	err = rdb.Decode(f, &decoder{})
	maybeFatal(err)
}

func parseArgTypes() {
	if printTypes == "" {
		printTypes = strings.Join(typesAll, ",")
	}
	arr := strings.Split(printTypes, ",")
	for _, v := range arr {
		v = strings.TrimSpace(v)
		if _, has := allKeyTypes[v]; has {
			kTypePrint[v] = true
		}
	}
}
