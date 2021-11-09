package uuid

import (
	"core/src/lib/random"
	"fmt"
)

func NewUUID(prefix string) string {
	hash := random.Snowflake().Base58()
	if prefix == ""{
		return hash
	}
	return fmt.Sprintf("%s%s",prefix, hash)
}