# Today I learned

Briefing some of the areas I've covered, lessons learned as part of this activity.
## Concepts
  - Channels
  - Go routines
  - Reader
  -  Timer and Ticker 
## Packages
   - io
   - os
   - csv
   - time
#### # Channels

Definition of channels in the Go Tour is as follows:

> Channels are a typed conduit through which you can send and receive values
> with the channel operator, <-. 
>By default, sends and receives block until the other side is ready.This allows goroutines to synchronize without explicit locks or condition variables.

#####  Buffered Channels
Definition of buffered channels in the Go Tour is as follows:

>Channels can be buffered. Provide the buffer length as the second argument to make to initialize a buffered channel:
>ch := make(chan int, 100)
>Sends to a buffered channel block only when the buffer is full. Receives block when the buffer is empty.

My questions:

- How channels and goroutines are related?Why do we need goroutines for using channel?
>Channels are the pipes that connect concurrent goroutines. You can send values into channels from one goroutine and receive those values into another goroutine.
> Refer: https://gobyexample.com/channels
- What is the deadlock of channel? When does this happen?
> A deadlock happens when a group of goroutines are waiting for each other and none of them is able to proceed.
> By default, sends and receives block until the other side is ready.

Note: 
>Main func can spawn other goroutines, but "main" itself is one groutine.

Deadlock: No reciever
#1
```
package main

func main() {
	messages := make(chan string)

	// Do nothing spawned goroutine
	go func() {}()

	// A goroutine ( main groutine ) trying to send message to channel
	// But no other goroutine runnning
	// And channel has no buffers
	// So it raises deadlock error
	messages <- "Sending" // fatal error: all goroutines are asleep - deadlock!
}
```
#2
```
package main

func main() {
        ch := make(chan int)
        ch <- 1 // stuck on the channel send operation that waits forever for someone to   // read the value. No other goroutine running, only "main" goroutine is running 
        // fatal error: all goroutines are asleep - deadlock!
        fmt.Println(<-ch) 
}
```
Deadlock: No sender

```
package main

func main() {
	messages := make(chan string)

	// Do nothing spawned goroutine
	go func() {}()

	// A goroutine ( main groutine ) trying to receive message from channel
	// But channel has no messages, it is empty.
	// And no other goroutine running. ( means no "Sender" exists )
	// So channel will be deadlocking
	<-messages // fatal error: all goroutines are asleep - deadlock!
}
```
- Why there is no deadlock in buffered channels even if we dont use goroutines?
> By default channels are unbuffered, meaning that they will only accept sends (chan <-) if there is a corresponding receive (<- chan) ready to receive the sent value. Buffered channels accept a limited number of values without a corresponding receiver for those values.

```
package main

import "fmt"

func main() {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	// Because this channel is buffered, we can send these values
	// into the channel without a corresponding concurrent receive (goroutine)
    // Later we can receive these two values as usual
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
\\ Output
\\ 1
\\ 2
```
#####  Range and Close

>A sender can close a channel to indicate that no more values will be sent. Receivers can test whether a channel has been closed by assigning a second parameter to the receive expression: after
v, ok := <-ch
>ok is false if there are no more values to receive and the channel is closed.
>The loop for i := range c receives values from the channel repeatedly until it is closed.
##### Select

>The select statement lets a goroutine wait on multiple channel operations.
>A select blocks until one of its cases can run, then it executes that case. It chooses one at random if multiple are ready.
##### Channel Synchronization
Question: Why does this program exit without starting the worker goroutine? How to make 
it work properly?
```
package main

import (
    "fmt"
    "time"
)

func worker(done chan bool) {
    fmt.Print("working...")
    fmt.Println("done")
    done <- true
}

func main() {
    done := make(chan bool, 1)
    go worker(done)
    fmt.Println("Exiting from main()")
    // Output
    // Exiting from main()
}
```
>In the above case, the main goroutine spawns another goroutine of worker function. Hence when we execute the above program, there are two goroutines running concurrently.Goroutines are scheduled cooperatively. Hence when the main goroutine starts executing, go scheduler do not pass control to the worker goroutine until the main goroutine does not execute completely. Unfortunately, when the main goroutine is done with execution, the program terminates immediately and scheduler did not get time to schedule worker goroutine.

Solution:
```
package main

import (
    "fmt"
    "time"
)

func worker(done chan bool) {
    fmt.Print("working...")
    fmt.Println("done")
    done <- true  
}

func main() {

    done := make(chan bool, 1)
    go worker(done)
    // time.Sleep(10 * time.Millisecond)
    // Either use time.Sleep or <-done
    <-done
     fmt.Println("Exiting from main()")
}
```
Solution:
- Using time.Sleep()
Before main goroutine pass control to the last line of code, we pass control to worker goroutine using time.Sleep(10 * time.Millisecond) call. In this case, the main goroutine sleeps for 10 milli-seconds and won’t be scheduled again for another 10 milliseconds.
- <-done: Channel Synchronisation
We can use channels to synchronize execution across goroutines. Using a blocking receive(<-done) will wait for the worker goroutine to finish.
#### # Goroutines
>A goroutine is a lightweight thread managed by the Go runtime.
>go f(x, y, z)
starts a new goroutine running
f(x, y, z)
The evaluation of f, x, y, and z happens in the current goroutine and the execution of f happens in the new goroutine.
Goroutines run in the same address space, so access to shared memory must be synchronized

Suppose we have a function call f(s). Here’s how we’d call that in the usual way, running it synchronously.
``` 
f("direct")
```
To invoke this function in a goroutine, use go f(s). 
```
go f("from goroutine")
```
This new goroutine will execute concurrently with the calling one.

You can also start a goroutine for an anonymous function call.
```
go func(msg string) {
        fmt.Println(msg)
    }("going")
```

Two important properties of goroutine:

- When a new Goroutine is started, the goroutine call returns immediately. Unlike functions, the control does not wait for the Goroutine to finish executing. The control returns immediately to the next line of code after the Goroutine call and any return values from the Goroutine are ignored.
- The main Goroutine should be running for any other Goroutines to run. If the main Goroutine terminates then the program will be terminated and no other Goroutine will run.

#### Reader

>Reader is any type that implements the `Read` method.The io.Reader interface represents an entity from which you can read a stream of bytes.
>A reader, represented by interface io.Reader, reads data from some source into a transfer buffer where it can be streamed and consumed

```
type Reader interface {
  Read(p []byte) (n int, err error)
}
```
Things to note:

- Read reads up to len(p) bytes into p and returns the number of bytes read – it returns an io.EOF error when the stream ends.
- After a Read() call, n may be less then len(p)
- Upon error, Read() may still return n bytes in buffer p
- If some data is available but not len(p) bytes, Read conventionally returns what is available instead of waiting for more.
- A call to Read() that returns n=0 and err=nil does not mean EOF as the next call to Read() may return more data.
- When a Read() exhausts available data, a reader may return a non-zero n and err=io.EOF. However, depending on implementation, a reader may choose to return a non-zero n and err = nil at the end of stream. In that case, any subsequent reads must return n=0, err=io.EOF.

Question:

- What is the use of Reader?
> You can read from a reader directly (this turns out to be the least useful use case):
```
p := make([]byte, 256)
n, err := r.Read(p)
```
Example: 

 Method Read is designed to be called within a loop where, with each iteration, it reads a chunk of data from the source and places it into buffer p. This loop will continue until the method returns an io.EOF error.
```
package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	r := strings.NewReader("abcde")

	buf := make([]byte, 4)
	for {
		n, err := r.Read(buf)
		fmt.Println(n, err, string(buf[:n]))
		if err == io.EOF {
			break
		}
	}
}
```
Output:
```
4 <nil> abcd
1 <nil> e
0 EOF 
```
Example that shows the contents in the transfer buffer:
```
package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	r := strings.NewReader("abcdef")

	buf := make([]byte, 4)

	for {
		n, err := r.Read(buf)
		fmt.Println(n, err, string(buf))
		if err != nil {
			if err == io.EOF {
				fmt.Println(string(buf))
				break
			}
			fmt.Println(err)
			return
		}
	}
}
```
Output:

```
4 <nil> abcd
2 <nil> efcd
0 EOF efcd
efcd
```

Use io.ReadFull to read exactly len(buf) bytes into buf:

```
package main

import (
	"fmt"
	"io"
	"log"
	"strings"
)

func main() {
	r := strings.NewReader("abcde")

	buf := make([]byte, 4)
	if _, err := io.ReadFull(r, buf); err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf))

	if _, err := io.ReadFull(r, buf); err != nil {
		fmt.Println(err)
	}
}
```
Output:
```
abcd
unexpected EOF
```
Use ioutil.ReadAll to read everything:

```
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	r := strings.NewReader("abcde")

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf))
}
```

Output:
```
abcde
```
Some examples of readers:
- If you open a file for reading, the object returned is an `os.File`, which is a Reader (it implements the Read method):
```
var r io.Reader
var err error
r, err = os.Open("file.txt")
```
- You can also make a Reader from a normal string using `strings.NewReader`:
```
var r io.Reader
r = strings.NewReader("Read will return these bytes")
```
- The body data from an `http.Request` is a Reader:
```
var r io.Reader
r = request.Body
```
- A bytes.Buffer is a Reader:
```
var r io.Reader
var buf bytes.Buffer
r = &buf
```
#### # Timer and Ticker
##### Timer

>Timers represent a single event in the future. You tell the timer how long you want to wait, and it provides a channel that will be notified at that time.

When the Timer expires, the current time will be sent on C, unless the Timer was created by AfterFunc.
```
type Timer struct {
    C <-chan Time
    // contains filtered or unexported fields
}
```

```
package main

import (
    "fmt"
    "time"
)

func main() {

    timer1 := time.NewTimer(2 * time.Second)

    <-timer1.C
    fmt.Println("Timer 1 fired")

    timer2 := time.NewTimer(time.Second)
    go func() {
        <-timer2.C
        fmt.Println("Timer 2 fired")
    }()
    stop2 := timer2.Stop()
    if stop2 {
        fmt.Println("Timer 2 stopped")
    }
    // Give the `timer2` enough time to fire, if it ever
	// was going to, to show it is in fact stopped.
    time.Sleep(2 * time.Second)
}
```
Output:
```
Timer 1 fired
Timer 2 stopped
```
The <-timer1.C blocks on the timer’s channel C until it sends a value indicating that the timer fired. If a program has already received a value from t.C, the timer is known to have expired and the channel drained.

Question:
- If timer is for waiting, why can't we just use time.Sleep()?
>One reason a timer may be useful is that you can Stop the timer before it fires.You can Reset() the timer if needed. 

##### Ticker

>Timers are for when you want to do something once in the future - tickers are for when you want to do something repeatedly at regular intervals
>Tickers can be stopped like timers. Once a ticker is stopped it won’t receive any more values on its channel.

Example of a ticker that ticks periodically until we stop it:
```
package main

import (
    "fmt"
    "time"
)

func main() {

    ticker := time.NewTicker(500 * time.Millisecond)
    done := make(chan bool)

    go func() {
        for {
            select {
            case <-done:
                return
            case t := <-ticker.C:
                fmt.Println("Tick at", t)
            }
        }
    }()

    time.Sleep(1600 * time.Millisecond)
    ticker.Stop()
    done <- true
    fmt.Println("Ticker stopped")
}
```
Output:
```
Tick at 2012-09-23 11:29:56.487625 -0700 PDT
Tick at 2012-09-23 11:29:56.988063 -0700 PDT
Tick at 2012-09-23 11:29:57.488076 -0700 PDT
Ticker stopped
```
Question:
- What is the difference between Timer and Ticker?
> In case of Timer, if a value from t.C is recieved, the timer is known to have expired and the channel is drained.
> In case of Ticker, even if a value from t.C is recived, the ticker is never expired until we stop it using ticker.Stop()

Timer Example:
```
package main

import (
	"fmt"
	"time"
)

func main() {

	timer := time.NewTimer(2 * time.Second)

	for i := 0; i < 3; i++ {
		t := <-timer.C
		fmt.Println(t)
		fmt.Printf("Fired %d time(s)\n", i+1)
	}

}
```

Output:
```
2009-11-10 23:00:02 +0000 UTC m=+2.000000001
Fired 1 time(s)
fatal error: all goroutines are asleep - deadlock!
```

Ticker Example:

```
package main

import (
	"fmt"
	"time"
)

func main() {

	ticker := time.NewTicker(2 * time.Second)

	for i := 0; i < 3; i++ {
		t := <-ticker.C
		fmt.Println(t)
		fmt.Printf("Fired %d time(s)\n", i+1)
	}

}
```
Output:
```
2009-11-10 23:00:02 +0000 UTC m=+2.000000001
Fired 1 time(s)
2009-11-10 23:00:04 +0000 UTC m=+4.000000001
Fired 2 time(s)
2009-11-10 23:00:06 +0000 UTC m=+6.000000001
Fired 3 time(s)
```