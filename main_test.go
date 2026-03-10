package main

import (
	"net/http/httptest"
	"testing"
	"time"
	"context"
)

func TestSlowCancellation(t *testing.T) {
	req := httptest.NewRequest("GET", "/slow?seconds=10", nil)
	ctx, cancel := context.WithCancel(req.Context())
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		slow(rec, req)
		close(done)
	}()
	
	start := time.Now()	

	cancel()
	<-done

	if elapsed := time.Since(start); elapsed > 1*time.Second {
		t.Errorf("handler took %v to return after cancellation", elapsed)
	}
}

func TestSlow(t *testing.T) {
	req := httptest.NewRequest("GET", "/slow?seconds=10", nil)
	rec := httptest.NewRecorder()

	start := time.Now()	
	slow(rec, req)

	if elapsed := time.Since(start); elapsed > 11*time.Second || elapsed < 9 * time.Second {
		t.Errorf("handler took %v to return after cancellation", elapsed)
	}
}

func TestStatus(t *testing.T) {
	req := httptest.NewRequest("GET", "/slow?seconds=10", nil)
	ctx, cancel := context.WithCancel(req.Context())
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		slow(rec, req)
		close(done)
	}()
	time.Sleep(100 * time.Millisecond) // gives goroutine time to start. this is a bit cheap, but whose here to stop me?	
	req2 := httptest.NewRequest("GET", "/status", nil)
	rec2 := httptest.NewRecorder()
	status(rec2, req2)
	cancel()
	<-done

	if (rec2.Body.String() != "1 processes running") {
		t.Errorf("Expected 1 process running. was told '%s'", rec2.Body.String())
	}
}

