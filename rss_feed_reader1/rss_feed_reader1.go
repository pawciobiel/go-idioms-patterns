package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Item struct {
	Title, Channel, GUID string // a subset of RSS fields
}

type Fetcher interface {
	Fetch() (items []Item, next time.Time, err error)
}

// fetches Items from domain
func Fetch(domain string) Fetcher {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}

type Subscription interface {
	Updates() <-chan Item // stream of Items
	Close() error         // shuts down the stream
}

// sub implements the Subscription interface.
type sub struct {
	fetcher Fetcher   // fetches items
	updates chan Item // delivers items to the user
	closed  bool
	err     error
}

// converts Fetches to a stream
func Subscribe(fetcher Fetcher) Subscription {
	s := &sub{
		fetcher: fetcher,
		updates: make(chan Item), // for Updates
	}
	go s.loop()
	return s
}

// loop fetches items using s.fetcher and sends them
// on s.updates.  loop exits when s.Close is called.
func (s *sub) loop() {
	for {
		if s.closed { // date race!
			close(s.updates)
			return
		}
		items, next, err := s.fetcher.Fetch()
		if err != nil {
			s.err = err                  // data race!
			time.Sleep(10 * time.Second) // sleeping
			continue
		}
		for _, item := range items {
			s.updates <- item // block!
		}
		if now := time.Now(); next.After(now) {
			time.Sleep(next.Sub(now)) // sleeping
		}
	}
}
func (s *sub) Updates() <-chan Item {
	return s.updates
}

func (s *sub) Close() error {
	s.closed = true // data race!
	return s.err    // data race!
}

// merges several streams
func Merge(subs ...Subscription) Subscription {

}

func main() {
	// Subscribe to some feeds, and create a merged update stream.
	merged := Merge(
		Subscribe(Fetch("blog.golang.org")),
		Subscribe(Fetch("googleblog.blogspot.com")),
		Subscribe(Fetch("googledevelopers.blogspot.com")))

	// Close the subscriptions after some time.
	time.AfterFunc(3*time.Second, func() {
		fmt.Println("closed:", merged.Close())
	})

	// Print the stream.
	for it := range merged.Updates() {
		fmt.Println(it.Channel, it.Title)
	}

	panic("show me the stacks")
}
