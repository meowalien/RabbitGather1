package errs

import (
	"core/src/lib/runtime"
	"fmt"
)


func WithLine(err error,msg ...interface{} )error {
	m := ""
	if msg != nil && len(msg)>=1{
		m = fmt.Sprintf(":%s",fmt.Sprint(msg...))
	}
	return fmt.Errorf("%s > %w%s",runtime.CallerFileAndLine(1) , err , m)
}
