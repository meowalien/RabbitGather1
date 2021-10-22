package random

import (
	math2 "core/sec/lib/math"
	"github.com/bwmarrin/snowflake"
	"math"
	"math/rand"
	"time"
)

var node *snowflake.Node

func init() {
	rand.Seed(time.Now().UnixNano())
}

func init() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err.Error())
	}
}

func Snowflake() snowflake.ID {
	return node.Generate()
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func GetSnowflakeIntWithLength(lg int64) int64 {
	return math2.CutIntMax(Snowflake().Int64(), lg)
}

func GetRandomInt(min int, max int) int {
	if min < 0 {
		panic("min < 0")
	}
	return rand.Intn(max-min) + min //CutIntBetween(Snowflake().Int64(), int64(math.Log10(float64(min)))+1, int64(math.Log10(float64(max)))+1)
}

func RandomInLength(i int) int {
	return rand.Intn(int(math.Floor(math.Pow(10.0, float64(i)))))
}

func randFloats(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}
func RandomByteArray(i int) []byte {
	token := make([]byte, i)
	rand.Read(token)
	return token
}
