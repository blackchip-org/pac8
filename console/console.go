package console

import (
	"fmt"
	"runtime"
	"sync"
)

var At = false
var mutex sync.Mutex

func Print(v ...interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	prefix()
	fmt.Print(v...)
}

func Printf(format string, a ...interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	prefix()
	fmt.Printf(format, a...)
}

func Println(v ...interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	prefix()
	fmt.Println(v...)
}

func prefix() {
	if !At {
		fmt.Print("(console)\t")
		return
	}
	_, file, line, _ := runtime.Caller(2)
	fmt.Printf("(console:%v:%v) ", file, line)
}
