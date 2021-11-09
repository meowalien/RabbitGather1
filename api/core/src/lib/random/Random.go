package random

import (
	"fmt"
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


func RangeInt(low, hi int) int {
	return low + rand.Intn(hi-low)
}

func MaxInt0Fill(length int ) string  {
	return fmt.Sprintf(fmt.Sprintf("%%0%dd",length),RangeInt(0,int(math.Pow(10, float64(length)))))
}

