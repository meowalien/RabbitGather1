package util

import (
	"github.com/bwmarrin/snowflake"
	"math"
	"math/rand"
	"time"
	"unsafe"
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
	return CutIntMax(Snowflake().Int64(), lg)
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

//func NewVerificationCodeWithLength(d int) string {
//	return fmt.Sprintf(fmt.Sprintf("%%0%dd", d), RandomInLength(d))
//}

func randFloats(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

var src = rand.NewSource(time.Now().UnixNano())

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
const letterBytes = "ABCDEFGHIJKLM_NOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"

func randStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func RandomString(i int) string {
	return randStringBytesMaskImprSrcUnsafe(i)
}

func RandomByteArray(i int) []byte {
	token := make([]byte, i)
	rand.Read(token)
	return token
}
