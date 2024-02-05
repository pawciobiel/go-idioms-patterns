# Go (Golang) idioms, code patterns, talks
Links to various information about the language in one place...

[Google I/O 2012 Go Concurrency Patterns - Rob Pike](https://www.youtube.com/watch?v=f6kdp27TYZs)
[source](https://go.dev/talks/2012/concurrency/support/)
1. [Using channels](#using-channels)
1. [Generator: function that returns a channel](#generator-function-that-returns-a-channel)
1. [Channels as handle on a service](#channels-as-handle-on-a-service)
1. [Multiplexing](#multiplexing)
1. [Restoring sequencing](#sequencing)
1. [Fan-in using select](#fanin_select)
1. [Timeout using select](#timeout_select)
1. [Timeout for whole conversation using select](#timeout_whole_select)
1. [Quit channel](#quit_channel)
1. [Receive on quit channel](#receive_on_quit_channel)
1. [Daisy-chain](#daisy_chain)
1. [Google Search](#google_search)

[GopherCon 2018: Bryan C. Mills - Rethinking Classical Concurrency Patterns](https://www.youtube.com/watch?v=5zXAHh5tJqQ)

[Google I/O 2013 - Advanced Go Concurrency Patterns](https://www.youtube.com/watch?v=QDDwwePbDtw)
[https://go.dev/talks/2013/advconc/](https://go.dev/talks/2013/advconc/)
1. [ping-pong](#ping_pong)


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



## Using channels

When the main function executes `<-c`, it will wait for a value to be sent.

Similarly, when the boring function executes `c <- value`, it waits for a
receiver to be ready.

A sender and receiver must both be ready to play their part in the communication.
Otherwise we wait until they are.

Thus channels both communicate and synchronize.

[using_channels/using_channels.go](using_channels/using_channels.go)


## Generator: function that returns a channel

Channels are first-class values, just like strings or integers.

[generator/generator.go](generator/generator.go)


## Channels as handle on a service
Our boring function returns a channel that lets us communicate with the boring service it provides.

We can have more instances of the service.

[hanle_on_a_service/hanle_on_a_service.go](hanle_on_a_service/hanle_on_a_service.go)


## Multiplexing
These programs make Joe and Ann count in lockstep.
We can instead use a fan-in function to let whosoever is ready talk.
```go:multiplexing/multiplexing.go

```


## Restoring sequencing
Send a channel on a channel, making goroutine wait its turn.
Receive all messages, then enable them again by sending on a private channel.
First we define a message type that contains a channel for the reply.
```go:sequencing/sequencing.go

```


## Fan-in using select
Rewrite our original fanIn function. Only one goroutine is needed. New:
```go:fanin_select/fanin_select.go

```


## Timeout using select
The time.After function returns a channel that blocks for the specified duration.
After the interval, the channel delivers the current time, once.
```go:timeout_select/timeout_select.go

```


## Timeout for whole conversation using select
Create the timer once, outside the loop, to time out the entire conversation.
(In the previous program, we had a timeout for each message.)
```go:timeout_whole_select/timeout_whole_select.go

```


## Quit channel
We can turn this around and tel Joe to stop when we're tired of listening to him.
```go:quit_channel/quit_channel.go

```


## Receive on quit channel
How do we know it's finished? Wait for it to tell us it's done: receive on the quite channel
```receive_on_quit_channel/receive_on_quit_channel.go

```


## Daisy-chain
100'000 goroutines
```daisy_chain/daisy_chain.go

```

## Google Search: A fake framework
We can simulate the search function, much as we simulated conversation before.
```google1.0/google1.0.go
```

Run the Web, Image and Video searches concurrently, and wait for all results.
No locks. No condition variables. No callbacks.
```google2.0/google2.0.go
```

Don't wait for slow servers. No locks. No condition variables. No callbacks.
```google2.1/google2.1.go
```

Avoid timeout
Q: How do we avoid discarding results from slow servers?
A: Replicate the servers. Send requests to multiple replicas, and use the first response.
```google2.2/google2.2.go
```

Replicas
Return the fastest from replicas.
```google3.0/google3.0.go
```

## ping-pong
```ping_pong/ping_pong.go
```

## RSS feed reader
Naive and buggy RSS reader
```rss_feed_reader1/rss_feed_reader1.go
```
