package entities

import (
	"log"
	"runtime"
)

func GoSafely(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				stack := make([]byte, 1024*8)
				stack = stack[:runtime.Stack(stack, false)]
				f := "PANIC: %s\n%s\n"
				log.Printf(f, r, stack)
			}
		}()
		fn()
	}()
}
