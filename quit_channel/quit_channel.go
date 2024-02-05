package main

import (
	"fmt"
	"math/rand"
	"time"
)

func boring(msg string, quit chan bool) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			select {
			case c <- fmt.Sprintf("%s %d", msg, i):
				// do nothing
			case <-quit:
				fmt.Println("ok bye!")
				return
			}
		}
	}()
	return c // Return the channel to the caller.
}

func main() {
	quit := make(chan bool)
	c := boring("Joe", quit)
	doneOn := rand.Intn(10)
	fmt.Printf("Done on %d\n", doneOn)
	for i := doneOn; i > 0; i-- {
		fmt.Println(<-c)
	}
	quit <- true
}
