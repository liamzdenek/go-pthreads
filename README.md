go-pthreads
===========

This is a binding of C's pthreads to Google Go. **This library is not a
replacement for goroutines.** This library is designed to help bind C libraries
with blocking function calls to Go in a go-friendly manner. If this is not your
use case, this library probably won't help you.

Use Case
--------

If a goroutine exists that calls a function that will block potentially
indefinitely, that goroutine cannot be stopped until the blocking function
returns and the goroutine checks an "exit" channel, and exits of its own will.

In every day go programming, this condition should not exist, as all inter-
thread communication, as well as reading and writing, should be done using
channels. However, in C, many libraries (and, in my specific case, networking
libraries) implement blocking functions (recv), and the mixture of a blocking
function and a multi-channel select caused many implementation problems and
"hacks" to get the "blocking" function to return periodically, so the exit
channel could be checked, and the routine could exit if it had to.

Example
-------

```go
package main;

import (
	"github.com/liamzdenek/go-pthreads"
	"fmt"
)

func main() {
	thread := pthread.Create(func() {
		// we're within the pthread
		counter := 1;
		for {
			// time to make a blocking function call. The library includes
			// pthread.sleep for demo purposes, but this will work with any
			// library that causes IO wait
		
			// an example using github.com/alecthomas/gozmq can be found in
			// Thread_test.go

			fmt.Printf("Hello, %d\n", counter)
			counter++
			pthread.Sleep(1) // seconds
		}
	})

	// within the main goroutine
	pthread.Sleep(3)

	fmt.Printf("Killing thread\n");
	thread.Kill()

	pthread.Sleep(3);
}
```

Output:
```
$ time go run test.go
Hello, 1
Hello, 2
Hello, 3
Killing thread
go run test.go  0.28s user 0.09s system 5% cpu 6.421 total
```

Pros/Cons
---------

Pros:

* Provides a mechanism to kill a blocked thread (thread.Kill())
* Provides thread status without any logic in the child (thread.Running())

Cons:

* Does not implement pthread_cleanup_push/pop
* Runs in a dedicated thread (most of the time; sometimes this is a pro)
* Does not integrate with the go scheduler (as a consequence of the new thread)
* Harder to debug (crashes in C code don't produce stack traces)

