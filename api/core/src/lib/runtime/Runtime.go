package runtime

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// 取得呼叫的文件與行號
func CallerFileAndLine(deap int ) string {
	_ ,file , line,ok := runtime.Caller(deap + 1)
	if !ok{
		return "[fail to get caller]"
	}
	file = filepath.Base(file)
	return fmt.Sprintf("%s:%d",file,line)
}