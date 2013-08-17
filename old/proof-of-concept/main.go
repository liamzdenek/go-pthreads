package main

/*
#include <pthread.h>

extern void callMe();

static void doCallMe(pthread_t* pid, pthread_attr_t* attr) {
	pthread_create(pid,attr,(void*)callMe,NULL);
}
*/
import "C"
import "fmt"
import "time"

//export callMe
func callMe() {
	time.Sleep(time.Second);
	fmt.Printf("In the C thread\n");
}

func main() {
	var attr C.pthread_attr_t
	var pid C.pthread_t

	C.pthread_attr_init(&attr)
	C.doCallMe(&pid, &attr);

	time.Sleep(time.Second*2)
	fmt.Printf("Bye\n")
}

