package main

import (
	"fmt"
	"math/rand"
	"time"
)

func cleanup() {
	fmt.Println("cleanup()")
}

func boring(msg string, quit chan string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			select {
			case c <- fmt.Sprintf("%s %d", msg, i):
				// do nothing
			case <-quit:
				cleanup()
				quit <- "See you!"
				return
			}
		}
	}()
	return c
}

func main() {
	quit := make(chan string)
	c := boring("Joe", quit)
	doneOn := rand.Intn(10)
	fmt.Printf("Finish on %d\n", doneOn)
	for i := doneOn; i > 0; i-- {
		fmt.Println(<-c)
	}
	quit <- "Bye!"
	fmt.Printf("Joe says: %q\n", <-quit)
}
