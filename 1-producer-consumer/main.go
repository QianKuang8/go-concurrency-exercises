//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream, tweets chan *Tweet) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(tweets)
			return
		}

		tweets <- tweet
	}
}

func consumer(wg *sync.WaitGroup, tweets chan *Tweet) {
	defer wg.Done()
	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	tweets := make(chan *Tweet, 100)

	// Producer
	go producer(stream, tweets)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	// Consumer
	go consumer(wg, tweets)
	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
