package usecases

import (
	"log"
	"runtime"
)

func CheckError(err error) bool {
	stack := make([]byte, 1024*8)
	stack = stack[:runtime.Stack(stack, false)]

	defer func() {
		if r := recover(); r != nil {
			f := "PANIC: %s\n%s\n"
			log.Printf(f, r, stack)
		}
	}()
	if err != nil {
		f := "ERROR : %s\n%s\n"
		log.Printf(f, err.Error(), stack)
		return true
	}

	return false
}

func PanicAndCheckError(err error) {
	stack := make([]byte, 1024*8)
	stack = stack[:runtime.Stack(stack, false)]

	defer func() {
		if r := recover(); r != nil {
			f := "PANIC: %s\n%s\n"
			log.Printf(f, r, stack)
		}
	}()
	if err != nil {
		panic(err)
	}
}
