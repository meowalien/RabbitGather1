package uuid

import (
	"core/sec/lib/random"
	"fmt"
)

var uuidPrefixMap = map[string]struct{}{}
func NewUUIDWithPrefix(prefix string)string {
	if _,ok:= uuidPrefixMap[prefix] ; ok{
		panic("duplicate prefix UUID: "+prefix)
	}
	uuidPrefixMap[prefix] = struct{}{}
	return fmt.Sprintf("%s_%s" , prefix, random.Snowflake().Base58())
}