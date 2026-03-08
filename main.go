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
		// we'll talk later
	}
	values := r.URL.Query()
	seconds := values.Get("seconds")
	counter_lock.Lock()
	counter += 1;
	counter_lock.Unlock()
	iseconds, _ := strconv.Atoi(seconds) // iseconds == 0 on error, which is fine.

	time.Sleep(time.Duration(iseconds) * time.Second)
	counter_lock.Lock()
	counter -= 1
	counter_lock.Unlock()

	fmt.Fprintf(w, "Waited %v seconds", iseconds);	
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v processes running", counter);	
}

func main() {
	http.HandleFunc("/slow", slow)
	http.HandleFunc("/status", status)


	log.Fatal(http.ListenAndServe(":8080", nil))
}
