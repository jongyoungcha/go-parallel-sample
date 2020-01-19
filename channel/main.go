package main

import (
	"fmt"
	"sync"
	"strconv"
	"os"
	"os/signal"
	"syscall"
)

type message struct {
	worker string
	data int
}


func main() {
	queue := make(chan message, 100)
	queueRets := make(chan int, 100)
	trigger := make(chan interface{})
	done := make(chan interface{})	
	
	wg := &sync.WaitGroup{}
	
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	
	// producer
	producer := func(wg *sync.WaitGroup, workerName string) {
		defer wg.Done()
		<-trigger
		for i := 1; i <= 10000; i++ {
			queue <- message{
				worker: workerName,
				data: i,
			}
		}
		fmt.Printf("Done (producer : %s)\n", workerName)
	}
	
	// consumer
	consumer := func(wg *sync.WaitGroup, workerName string) {
		defer wg.Done()
		<-trigger

		count:=0
		
		for {
			select {
			case msg := <- queue:
				fmt.Printf("received (%s) : %s:%d\n", workerName, msg.worker, msg.data)
				count++
				
			case <-done:
				queueRets <- count
				return
			}
		}
		
		fmt.Println("Consumer done!")
	}
	
	// make producers
	for i:=0; i<10; i++ {
		wg.Add(1)
		go producer(wg, "p"+strconv.Itoa(i))
	}

	// make consumers
	for i:=0; i<10; i++ {
		wg.Add(1)
		go consumer(wg, "c"+strconv.Itoa(i))
	}
	
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println("Received Signal: ", sig)
		close(done)
	}()

	// Start workers
	close(trigger)
	
	// Wait groups
	wg.Wait()
	close(queue)
	
	total:=0
	for i:=0; i<10; i++ {
		sum:= <-queueRets
		total += sum
		fmt.Printf("sum of consumer(%d):%d\n", i, sum)
	}
	
	fmt.Println("total:", total)
	fmt.Println("Done...")
}
