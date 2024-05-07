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
	"log"
	"sync"
	"time"
)

var mu sync.Mutex

func producer(stream Stream, tweets chan Tweet) {
	defer close(tweets)
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			log.Printf("erro %v", err)
			break
		}

		tweets <- *tweet
	}
}

func consumer(tweets <-chan Tweet) {
	for {
		tweet, ok := <-tweets
		if !ok {
			return
		}
		if tweet.IsTalkingAboutGo() {
			fmt.Println(tweet.Username, "\ttweets about golang")
		} else {
			fmt.Println(tweet.Username, "\tdoes not tweet about golang")
		}
		//		time.Sleep(400 * time.Millisecond)
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()
	tweets := make(chan Tweet)
	var wg sync.WaitGroup

	wg.Add(2)
	// Producer
	go func() {
		defer wg.Done()
		producer(stream, tweets)
	}()

	go func() {
		defer wg.Done()
		consumer(tweets)
	}()

	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
