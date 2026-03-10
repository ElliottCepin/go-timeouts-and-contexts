package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"log"
	"sync"
)

var counter_lock sync.Mutex
var counter int = 0


func slow(w http.ResponseWriter, r *http.Request) {
	if (r.Method != "GET"){
		return
	}
	ctx := r.Context()

	result := make(chan int, 1)
	
	go func() {
		values := r.URL.Query()
		seconds := values.Get("seconds")
		counter_lock.Lock()
		counter += 1;
		counter_lock.Unlock()
		iseconds, _ := strconv.Atoi(seconds) // iseconds == 0 on error, which is fine.
		select {
		case <-time.After(time.Duration(iseconds) * time.Second):
		case <-ctx.Done():
		}
		
		counter_lock.Lock()
		counter -= 1
		counter_lock.Unlock()
		result <- iseconds
	}()
	
	select {
		case res := <- result:
			fmt.Fprintf(w, "Waited %v seconds", res);	
		case <- ctx.Done():
			log.Print("Client timed out")
	}

}

func status(w http.ResponseWriter, r *http.Request) {
	counter_lock.Lock()
	fmt.Fprintf(w, "%v processes running", counter);
	counter_lock.Unlock()
}

func main() {
	srv := &http.Server{
		Addr:		":8080",
		ReadTimeout:	5 * time.Second,
		WriteTimeout:	30 * time.Second,
	}
	http.HandleFunc("/slow", slow)
	http.HandleFunc("/status", status)


	log.Fatal(srv.ListenAndServe())
}
