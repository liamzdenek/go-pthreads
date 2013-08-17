package pthread

import (
	"github.com/alecthomas/gozmq"
	"testing"
	"time"
)

func Test_ZMQ(t *testing.T) {
	// initialize some ZMQ crap
	failed := false
	context, _ := gozmq.NewContext()
	socket, _ := context.NewSocket(gozmq.REP)
	socket.Bind("tcp://127.0.0.1:9898")

	// create the thread
	thread := Create(func() {
		// this call should block forever, and thread.Kill() should stop it
		socket.Recv(0)

		// this should never happen
		failed = true
	})

	// wait around for a bit
	time.Sleep(time.Millisecond * 100)

	// make sure the thread started
	if !thread.Running() {
		t.Error("The thread was not running after starting it")
	}

	// stop the thread
	thread.Kill()

	// wait around a bit more
	time.Sleep(time.Millisecond * 100)

	// make sure the thread terminated
	if thread.Running() {
		t.Error("The thread was still running after the kill signal")
	}

	// finally, make sure that the thread didn't go to a point in execution that
	// it was not supposed to reach
	if failed {
		// this might be more revealing of a bug in zmq or pthreads, but what do
		// i know
		t.Error("The recv call did not block or the thread did not exit properly.")
	}

}

func Test_ManyThreads(t *testing.T) {
	literal := func() {
		for {
			Sleep(1); // sleep for 1s, forever
		}
	}

	// start up the threads
	thread1 := Create(literal);
	thread2 := Create(literal);

	// give them some time to spin up
	time.Sleep(time.Millisecond * 100);

	// ensure both started
	if !thread1.Running() || !thread2.Running() {
		t.Error("One or both of the threads failed to start properly");
	}

	// kill one of the threads
	thread1.Kill();

	// give it some time to clean up
	time.Sleep(time.Millisecond * 100);

	// ensure one has exited and the other has running
	if thread1.Running() {
		t.Error("Thread 1 has failed to stop");
	}
	if !thread2.Running() {
		t.Error("Thread 2 stopped when it shouldn't have");
	}

	// stop the second thread
	thread2.Kill();

	// wait for it to spin down
	time.Sleep(time.Millisecond * 100);

	// ensure that it has spun down
	if thread2.Running() {
		t.Error("Thread 2 has failed to stop");
	}
}
