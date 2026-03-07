package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"log"
)

func slow(w http.ResponseWriter, r *http.Request) {
	if (r.Method != "GET"){
		// we'll talk later
	}
	values := r.URL.Query()
	seconds := values.Get("seconds")
	iseconds, _ := strconv.Atoi(seconds) // iseconds == 0 on error, which is fine.
	time.Sleep(time.Duration(iseconds) * time.Second)
	fmt.Fprintf(w, "Waited %v seconds", iseconds);	
}

func main() {
	http.HandleFunc("/slow", slow)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
