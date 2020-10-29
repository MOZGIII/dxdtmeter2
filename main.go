package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type Counter struct {
	val uint64
}

func NewCounter() *Counter {
	return &Counter{val: 0}
}

func (c *Counter) Inc() {
	atomic.AddUint64(&c.val, 1)
}

func (c *Counter) Add(val uint64) {
	atomic.AddUint64(&c.val, val)
}

func (c *Counter) Get() uint64 {
	return atomic.LoadUint64(&c.val)
}

func (c *Counter) Reset() {
	atomic.StoreUint64(&c.val, 0)
}

func serve(addrEnvVarName string, hanlder http.Handler) {
	addr := os.Getenv(addrEnvVarName)
	if len(addr) == 0 {
		log.Fatalf("Unable to find env var %s", addrEnvVarName)
	}

	go func() {
		if err := http.ListenAndServe(addr, hanlder); err != nil {
			log.Fatalf("Unable to start the server: %v", err)
		}
	}()

	log.Printf("Listening at %s from %s", addr, addrEnvVarName)
}

func main() {
	counter := NewCounter()

	mode := os.Getenv("MODE")
	var incHandler http.Handler
	switch mode {
	case "ndjson", "jsonlines", "lines":
		incHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			buf := make([]byte, 1024)
			for {
				n, err := r.Body.Read(buf)
				if err != io.EOF || n == 0 {
					break
				}
				counter.Add(uint64(bytes.Count(buf[:n], []byte("\n"))))
			}
		})
	case "", "simple":
		incHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter.Inc()
		})
	default:
		log.Fatalf("Invalid mode %q", mode)
	}
	log.Printf("Mode: %s", mode)

	controlHandler := http.NewServeMux()
	controlHandler.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		counter.Reset()
	})
	controlHandler.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%d", counter.Get())
	})

	serve("ADDR", incHandler)
	serve("CONTROL_ADDR", controlHandler)

	log.Print("Booted")
	<-make(chan struct{})
}
