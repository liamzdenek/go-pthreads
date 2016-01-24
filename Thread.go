package pthread

/*
#include <pthread.h>
#include <signal.h>
#include <unistd.h>
#include <stdio.h>

extern void createThreadCallback();
static void sig_func(int sig);

static void createThread(pthread_t* pid) {
	pthread_create(pid, NULL, (void*)createThreadCallback, NULL);
}

static void sig_func(int sig)
{
	//printf("handling exit signal\n");
	signal(SIGSEGV,sig_func);
	pthread_exit(NULL);
}

static void register_sig_handler() {
	signal(SIGSEGV,sig_func);
}
*/
import "C"
import "unsafe"

type Thread uintptr
type ThreadCallback func(args interface{})

var create_callback chan ThreadCallback
var create_args chan interface{}

// var create_args     chan interface

func init() {
	C.register_sig_handler()
	create_callback = make(chan ThreadCallback, 1)
	create_args = make(chan interface{}, 1)
}

//export createThreadCallback
func createThreadCallback() {
	C.register_sig_handler()
	C.pthread_setcanceltype(C.PTHREAD_CANCEL_ASYNCHRONOUS, nil)
	callback := <-create_callback
	args := <-create_args
	callback(args)

}

// calls C's sleep function
func Sleep(seconds uint) {
	C.sleep(C.uint(seconds))
}

// initializes a thread using pthread_create
func Create(cb ThreadCallback, args interface{}) Thread {
	var pid C.pthread_t
	pidptr := &pid
	create_callback <- cb
	create_args <- args
	C.createThread(pidptr)

	return Thread(uintptr(unsafe.Pointer(&pid)))
}

// determines if the thread is running
func (t Thread) Running() bool {
	// magic number "3". oops
	// couldn't figure out the proper way to do this. probably because i suck
	// if someone knows the right way, pls submit a pull request
	return int(C.pthread_kill(t.c(), 0)) != 3
}

// signals the thread in question to terminate
func (t Thread) Kill() {
	C.pthread_kill(t.c(), C.SIGSEGV)
}

// helper function to convert the Thread object into a C.pthread_t object
func (t Thread) c() C.pthread_t {
	return *(*C.pthread_t)(unsafe.Pointer(t))
}
