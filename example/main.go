package main

import (
	_ "github.com/tam7t/sigprof"

	"fmt"
	"time"
)

func main() {
	messages := make(chan string)

	// consumer
	go func() {
		for {
			select {
			case m := <-messages:
				fmt.Println(m)
			}
		}
	}()

	// producer
	go func() {
		for {
			time.Sleep(1 * time.Second)
			messages <- `ping`
		}
	}()

	// block indefinately
	select {}
}
