package errs

import (
	"core/src/lib/runtime"
	"fmt"
)


func WithLine(err error , deap ...int )error {
	dp := 1
	if len(deap) == 1{
		dp += deap[0]
	}
	line := runtime.CallerFileAndLine(dp)
	return formatErrLine(err , line )
}

func formatErrLine ( err error , line string )error{
	return fmt.Errorf("%s - %w",line , err )
}

func Join(errs ...error )error {
	var finalErr error
	for _, err := range errs {
		if err == nil{
			continue
		}
		if finalErr == nil{
			finalErr = err
		}else{
			finalErr = fmt.Errorf("%w > %s" ,finalErr , err.Error() )
		}
	}
	return finalErr
}

type OneErrorReturnFuc func() error
type LogFunc func(...interface{})

func LogIfErr(f OneErrorReturnFuc , logger LogFunc)  {
	if e:= f() ; e!= nil{
		logger(e)
	}
}
