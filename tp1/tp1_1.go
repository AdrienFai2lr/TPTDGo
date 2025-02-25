package main

import (
	"fmt"
)

func main() {
	sync := make(chan bool) //synchro via canal

	go func() { //goroutine fct anonyme
		fmt.Println("Hello")
		//sync <- true
	}()
	fmt.Println("World")
	<-sync //dead lock
}
