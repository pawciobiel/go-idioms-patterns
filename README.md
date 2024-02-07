# Go (Golang) idioms, code patterns, talks
Links to various information about the language in one place...

[Google I/O 2012 Go Concurrency Patterns - Rob Pike](https://www.youtube.com/watch?v=f6kdp27TYZs)
[source](https://go.dev/talks/2012/concurrency/support/)
1. [Using channels](#using-channels)
1. [Generator: function that returns a channel](#generator-function-that-returns-a-channel)
1. [Channels as handle on a service](#channels-as-handle-on-a-service)
1. [Multiplexing](#multiplexing)
1. [Restoring sequencing](#restoring-sequencing)
1. [Fan-in using select](#fan-in-using-select)
1. [Timeout using select](#timeout-using-select)
1. [Timeout for whole conversation using select](#timeout-for-whole-conversation-using-select)
1. [Quit channel](#quit-channel)
1. [Receive on quit channel](#receive-on-quit-channel)
1. [Daisy-chain](#daisy-chain)
1. [Google Search](#google-search-a-fake-framework)

[GopherCon 2018: Bryan C. Mills - Rethinking Classical Concurrency Patterns](https://www.youtube.com/watch?v=5zXAHh5tJqQ)
1. [Future: API](#future-api)

[Google I/O 2013 - Advanced Go Concurrency Patterns](https://www.youtube.com/watch?v=QDDwwePbDtw)
[https://go.dev/talks/2013/advconc/](https://go.dev/talks/2013/advconc/)
1. [ping-pong](#ping-pong)
1. [RSS feed reader](#rss-feed-reader)

[Twelve Go Best Practices - Francesc Campoy](https://www.youtube.com/watch?v=8D3Vmm1BGoY)

[Golang UK Conference 2016 - Idiomatic Go Tricks - Mat Ryer](https://www.youtube.com/watch?v=yeetIgNeIkc)


Google golang talks:
1. [https://go.dev/talks/2011/](https://go.dev/talks/2011/)
1. [https://go.dev/talks/2012/](https://go.dev/talks/2012/)
1. [https://go.dev/talks/2013/](https://go.dev/talks/2013/)
1. [https://go.dev/talks/2014/](https://go.dev/talks/2014/)
1. [https://go.dev/talks/2015/](https://go.dev/talks/2015/)
1. [https://go.dev/talks/2016/](https://go.dev/talks/2016/)
1. [https://go.dev/talks/2017/](https://go.dev/talks/2017/)
1. [https://go.dev/talks/2019/](https://go.dev/talks/2019/)


## Google I/O 2012 Go Concurrency Patterns - Rob Pike
### Using channels

When the main function executes `<-c`, it will wait for a value to be sent.

Similarly, when the boring function executes `c <- value`, it waits for a
receiver to be ready.

A sender and receiver must both be ready to play their part in the communication.
Otherwise we wait until they are.

Thus channels both communicate and synchronize.

[using_channels/using_channels.go](using_channels/using_channels.go)


### Generator: function that returns a channel

Channels are first-class values, just like strings or integers.

[generator/generator.go](generator/generator.go)


### Channels as handle on a service
Our boring function returns a channel that lets us communicate with the boring service it provides.

We can have more instances of the service.

[hanle_on_a_service/hanle_on_a_service.go](hanle_on_a_service/hanle_on_a_service.go)


### Multiplexing
These programs make Joe and Ann count in lockstep.
We can instead use a fan-in function to let whosoever is ready talk.
[multiplexing/multiplexing.go](multiplexing/multiplexing.go)


### Restoring sequencing
Send a channel on a channel, making goroutine wait its turn.
Receive all messages, then enable them again by sending on a private channel.
First we define a message type that contains a channel for the reply.
[sequencing/sequencing.go](sequencing/sequencing.go)


### Fan-in using select
Rewrite our original fanIn function. Only one goroutine is needed. New:
[fanin_select/fanin_select.go](fanin_select/fanin_select.go)


### Timeout using select
The time.After function returns a channel that blocks for the specified duration.
After the interval, the channel delivers the current time, once.
[timeout_select/timeout_select.go](timeout_select/timeout_select.go)


### Timeout for whole conversation using select
Create the timer once, outside the loop, to time out the entire conversation.
(In the previous program, we had a timeout for each message.)
[timeout_whole_select/timeout_whole_select.go](timeout_whole_select/timeout_whole_select.go)


### Quit channel
We can turn this around and tel Joe to stop when we're tired of listening to him.
[quit_channel/quit_channel.go](quit_channel/quit_channel.go)


### Receive on quit channel
How do we know it's finished? Wait for it to tell us it's done: receive on the quite channel
[receive_on_quit_channel/receive_on_quit_channel.go](receive_on_quit_channel/receive_on_quit_channel.go)


### Daisy-chain
100'000 goroutines
[daisy_chain/daisy_chain.go](daisy_chain/daisy_chain.go)

### Google Search: A fake framework
We can simulate the search function, much as we simulated conversation before.
[google1.0/google1.0.go](google1.0/google1.0.go)

Run the Web, Image and Video searches concurrently, and wait for all results.
No locks. No condition variables. No callbacks.
[google2.0/google2.0.go](google2.0/google2.0.go)


Don't wait for slow servers. No locks. No condition variables. No callbacks.
[google2.1/google2.1.go](google2.1/google2.1.go)

Avoid timeout
Q: How do we avoid discarding results from slow servers?
A: Replicate the servers. Send requests to multiple replicas, and use the first response.
[google2.2/google2.2.go](google2.2/google2.2.go)

Replicas
Return the fastest from replicas.
[google3.0/google3.0.go](google3.0/google3.0.go)


## GopherCon 2018: Bryan C. Mills - Rethinking Classical Concurrency Patterns

### Future: API
```go
func Fetch(ctx context.Context, name string) <- chan Item {
    c := make(chan Item, 1)
    go func() {
        //[...]
        c <- item
        // If the item does not exist, 
        // Fetch closes the channel without sending.
    }()
    return c
}

a := Fetch(ctx, "a")
b := Fetch(ctx, "b")
consume(<-a, <-b)

```
* if we return without waiting for the future to complete, how long would they continue using resources?
* may we start fetches faster than we can retire them and run out of memory?
* will `Fetch()` keep using past in context after it returns?
* - if so what happens if we cancel it and then try to read from the channel?
  - will we receive a zero value? some other sentinel value? will we block?

### Producer-Consumer Queue: API
```go
// Glob finds all items with names matching pattern
// and sends them on the returned channel.
// It closes the channel when all items have been sent.
func Glob(ctx context.Context, pattern string) <-chan Item {
    go func() {
        defer close(c)
        for [...] {
            [...]
            c <- item
        }
    }()
    return c
}


for item := range Glob(ctx, "[ab]*") {
    [...]
}
```
* if we return without draining the channel from `Glob()` would we leak the goroutine that is sending to it?
* would `Glob()` keep using past in context as we iterate of the result?
* if so what happens if we cancel it? will we still get results?
  - when if ever will the channel be closed in that case?

### Caller-side concurrency: synchronous API
The caller can invoke synchronous functions concurrently,
and often won't need to use channels at all.
```go
func Fetch(ctx context.Context, name string) (Item error) {
    [...]
}

var a, b Item
g, ctx := errgroup.WithContext(ctx)
g.Go(func() (err error) {
    a, err = Fetch(ctx, "a")
    return err
})
g.Go(func() (err error) {
    b, err = Fetch(ctx, "b")
    return err
})
err := g.Wait()
[...]
consume(a, b)
```

### Internal, caller-side concurrency: channels
*Make concurrency an internal detail*
A synchronous API can have a concurrent implementation.


```go
func Glob(ctx context.Context, pattern string) ([]Item, error) {
    [...]  // Find matching names.
    c := make(chan Item)
    g, ctx := errgroup.WithContext(ctx)
    for _, name := range names {
        name := name
        g.Go(func() error {
            item, err := Fetch(ctx, name)
            if err == nil {
                c <- item
            }
            return err
        })
    }
    
    go func() {
        err = w.Wait()
        close(c)
    }()
    
    var items []Item
    for item := range c {
        items = append(items, item)
    }
    if err != nil {
        return nil, err
    }
    return items, nil
}
```


### Condition variable: setup
```go
type Queue struct {
    mu sync.Mutex
    items []Item
    itemAdded sync.Cond
}

func NewQueue() *Queue {
    q:= new(Queue)
    q.itemAdded.L = &q.mu
    return q
}
```
_A condition variable is associated with a `sync.Locker` (e.g., a Mutex)._

### Condition variable: wait and signal
```go

func (q *Queue) Get() Item {
    q.mu.Lock()
    defer q.mu.Unlock()
    for len(q.items) == 0 {
        q.itemAdded.Wait()
    }
    item := q.items[0]
    q.items = q.items[1:]
    return item
}

func (q *Queue) Put(item Item) {
    q.mu.Lock()
    defer q.mu.Unlock()
    q.items = append(q.items, item)
    q.itemAdded.Signal()
}
```
_Wait atomically unlocks the mutex and suspends the goroutine._
_Signal locks the mutex and wakes up the goroutine._

### Condition variable: broadcast
```go
type Queue struct {
    [...]
    closed bool
}

func (q *Queue) Close() {
    q.mu.Lock()
    defer q.mu.Unlock()
    q.closed = true
    q.cond.Broadcast()
}
```
_`Broadcast` usually communicates events that affects all waiters._


```go
func(q *Queue) GetMany(n int) []Item {
    q.mu.Lock()
    defer q.mu.Unlock()
    for len(q.items) < n {
        q.itemAdded.Wait()
    }
    item := q.items[:n:n]
    q.items = q.items[n:]
    return item
}

func (q *Queue) Put(item Item) {
    q.mu.Lock()
    defer q.mu.Unlock()
    q.items = append(q.items, item)
    q.itemAdded.Broadcast()
}

```

_Since we don't know **which** of `GetMany` calls may be ready after a `Put()`,
we can wake them **all** up and let them decide._


* Spurious wakeups - for events that aren't really global `Broadcast()` may wake up too many waiters, eg
 - one call to `Put()` wakes up all the `GetMany()` callers even though at most only one of them will be able to complete.
 - even `Singlal()` can result in spurious wakeups - it could wake up a caller that is not read instead of the one that it is
 - if it does that repeatedly it could strand items in the queue without the corresponding wake ups
 (avoiding spurious wakeups adds even more complexity and subtlety to the code)
* Forgotten signals
 - since the condition variable decouples the signal from the data it's easy to add some new code that updates the data
 and forget to signal the condition.
* Starvation - even if we don't forget a signal if the waiters are not uniform the pickier ones can starve
 - suppose we have one `GetMany(3000)` and once caller executing `GetMany(3)` in a tight loop - the two waiters will be 
 about equally likely to wake up but the `GetMany(3)` will be able to consume 3 items every 3 calls but `GetMany(3000)` won't
 have enough ready - the queue will remain drained and the larger call will block forever... we could add explicit way to
 avoid starvation but that again makes the code more complex.
* Unresponsive cancellation - the whole point of condition variables is to put goroutines to sleep while we wait for something to happen,
 but while we are waiting for that condition we may miss some other event that we ought to notice too - for example
 the caller may decide they don't want to wait that long and cancel the past in context expecting
 us to notice and return more or less immediately.
 **Unfortunately condition variables only lets us wait for events associated with their own mutex.** - so we can't select on
 a condition and a cancellation in the same time. Even if the caller cancels our call will block until the next
 time the condition is signaled.


### Condition variable: resource pool
```go

type Pool struct {
    mu sync.Mutex
    cond sync.Cond
    numConns, limit int
    idle []net.Conn
}

func NewPool(limit int) *Pool {
    p := &Pool{limit: limit}
    p.cond.L = &p.mu
    return p
}

func (p *Pool) Release(c net.Conn) {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.idle = append(p.idel, c)
    p.cond.Signal()
}

func (p *Pool) Hijack(c net.Conn) {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.numConns--
    p.cond.Signal()
}

func (p *Pool) Acquire(c net.Conn) net.Conn, error {
    p.mu.Lock()
    defer p.mu.Unlock()
    for len(p.idle) >= p.limit {
        p.cond.Wait()
    }
    
    if len(p.idle) >0 {
        c := p.idle[len(p.idle) - 1]
        p.idle = p.idle[:len(p.idle) - 1]
        return c, nil
    }
    
    c, err := dial()
    if err == nil {
        p.numConns++
    }
    return c, err

    
}
```
# Communication: resource pool

> A buffered channel can be used like a semaphore [...].
> The capacity of the channel buffer limits the number of simultaneous calls
> to process.
>
> Effective Go

```go
type Pool struct {
    sem chan token
    idle chan net.Conn
}

type token struct{}

func NewPool(limit int) * Pool {
    sem := make(chan token, limit)
    idle := make(chan net.Conn, limit)
    return &Pool{sem, idle}
}

func (p *Pool) Release(c net.Conn) {
    p.idle <= c
}

func (p *Pool) Hijack(c net.Conn) {
    <-p.sem
}

func (p *Pool) Acquire(ctx context.Context) (net.Conn, error) {
    select {
    case conn := <-p.idle:
        return conn, nil
    case p.sem <- token{}:
        conn, err := dial()
        if err != nil {
            <-p.sem
        }
        return conn, err
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

```

* channel operations combine synchronization, signaling, and data transfer.
* when we block on communicating, others can also communicate with us:
  for example, to cancel the call.

# communication: queue
```go
type Queue struct {
    items chan [Item]  // non-empty slices only
    empty chan bool    // holds true if the queue is empty
}

func NewQueue() *Queue {
    items := make(chan []Item, 1)
    empty := make(chan bool, 1)
    empty <- true
    return &Queue{items, empty}
}

func (q *Queue) Get(ctx context.Context) Item {
    var items []Item
    select {
        case <-ctx.Done():
            return 0, ctx.Err()
        case items = <-q.items:
    }
    
    item := items[0]
    if len(items) == 1 {
        q.empty <- true
    } else {
        q.items <- items[1:]
    }
    return item, nil
}

func (q *Queue) Put(item Item) {
    var items []Item
    select {
        case items = <-q.items:
        case <-q.empty:
    }
    items = append(items, item)
    q.items <- items
}
```

### specific communication: queue

```go
type waiter struct {
    n int
    c chan []Item
}

type state struct {
    items []Item
    wait []waiter
}

type Queue struct {
    s chan state
}

func NewQueue() *Queue {
    s := make(chan state, 1)
    s <- state{}
    return &Queue{s}
}

func (q *Queue) GetMany(n int) []Item{
    s := <-q.s
    if len(s.wait) == 0 && len(s.items) >= n {
        items := s.items[:n:n]
        s.items = s.items[n:]
        q.s <- s
        return items
    }
    c := make(chan []Item)
    s.wait = append(s.wait, waiter{n, c})
    q.s <- s
    return <-c
}

func (q *Queue) Put(item Item) {
    s := <-q.s
    s.items = append(s.items, item)
    for len(s.wait) > 0 {
        w := s.wait[0]
        if len(s.items) <- w.n {
            break
        }
        w.c <- s.items[:w.n:w.n]
        w.items = s.items[w.n:]
        s.wait = s.wait[1:]
    }
    q.s <- s
}

// TODO cancelation
```


### condition variable: repeating transition
```go
type Idler struct {
    mu syn.Mutex
    idle sync.cond
    busy bool
    idles int64
}

func (i *Idler) AwaitIdle() {
    i.mu.Lock()
    defer i.mu.Unlock()
    idles := i.idles
    for i.busy && idles == i.idles {
        i.idle.Wait()
    }
}

func (i *Idler) SetBusy(b bool) {
    i.mu.Lock()
    defer i.mu.Unlock()
    wasBusy := i.busy
    i.busy = b
    if wasBusy && !i.busy {
        i.idles++
        i.idle.Broadcast()
    }
}

func NewIdler() *Idler {
    i := new(Idler)
    i.idle.L = &i.mu
    return i
}
```

* we need to store state explicitly (One may think we need to store only current state - the busy boolean,
  but that turns out to be very subtle decision.) If `AwaitIdle` boot only until it saw a non busy state
  it would be possible to boot from busy to idle and back before `AwaitIdle` got a chance to check,
  and we would miss short idle events.
* go condition variables unlike pthread condition variables don't have spurious wakeups so in theory
  we could return from `AwaitIdle()` unconditionally after the first wait call.
* it's common for condition based code to overs-signal - for eg, to work around a non diagnosed deadlock - 
  so to avoid introducing subtle problems latter it's best to keep the code robust to spurious wakeups.
  Instead we could track cumulative counting events and wait until we either catch the idle events in progress
  or observe it's effect on a counter.


### communication: repeating transition
**We could avoid the double state transition entirely by communicating the transition instead of signaling it
 and we can plum in the cancelation to boot. We can broadcast a state transition by closing a channel**
- a state transition marks the completion of the previous state
- and closing a channel marks the completion of communication of that channel
```go
type Idler struct {
    next chan chan struct{}
}

func (i *Idler) AwaitIdle(ctx context.Context) error {
    idle := <- i.next
    i.next <- idle
    if idle != nil {
        select {
            case <-ctx.Done():
                return ctx.Err()
            case <-idle:
                // idle
        }
    }
    return nil
}

func (i *Idler) SetBusy(b bool) {
    idle := <- i.next
    if b && (idle == nil) {
        idle = make(chan struct{})
    } else if ~b && (idle != nil) {
        close(idle)  /// Idle now.
        idle = nil
    }
    i.next <- idle
}

func NewIdler() *Idler {
    next := make(chan [...], 1)
    next <- nil
    return &Idler{next}
}

```
_Allocate the channel ti be closed when the event starts,
or when the first waiter appear._

### worker pool

```go
work := make(chan Task)
for n := limit; n > 0; n-- {
    go func() {
        for task := range work {
            perform(task)
        }
    }()
}

for _, task := range hudgeSlice {
    work <- task  // sender blocks untill the worker is available to receive
}
```
* leaks workers forever!

```go
work := make(chan Task)
var wg sync.WaitGroup
for n := limit; n > 0; n-- {
    wg.Add(1)
    go func() {
       for task := range work {
           perform(task)
       }
       wg.Done()
    }()
}

for _, task := range hudgeSlice {
    work <- task
}
close(work)
wg.Wait()
```

### waitgroup: distributed (unlimited) work
**Start goroutine when you have concurrent work _to do now_.**
- and let them exit as soon as the work is done...

```go
var wg sync.WaitGroup
for _, task := range hudgeSlice {
    wg.Add(1)
    go func(task Task) {
        perform(task)
        wg.Done()
    }(task)
}
wg.Wait()
```
- but now we need to figure out how to limit the work again... 

### semaphore channel: limiting concurrency
 - semaphore channel: inverted worker pool
```go
sem := make(chan token, limit)
for _, task := range hudgeSlice {
    sem <- token{}
    go func(task Task) {
        perform(task)
        <-sem
    }(task)
}
for n := limit; n > 0; n-- {  // wait for the last tasks to finish
    sem <- token{}
}
```
* WaitGroup allows further add calls during wait while our `sem` does not.



## Google I/O 2013 - Advanced Go Concurrency Patterns

### ping-pong
[ping_pong/ping_pong.go](ping_pong/ping_pong.go)

### RSS feed reader
Naive and buggy RSS reader [rss_feed_reader1/rss_feed_reader1.go](rss_feed_reader1/rss_feed_reader1.go)
Fixed rss reader [realmain.go](https://go.dev/talks/2013/advconc/realmain/realmain.go)
