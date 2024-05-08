//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

var mu sync.Mutex

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	done := make(chan bool)
	timeUp := make(chan bool)
	var result bool

	go func() {
		go func() {
			for {
				mu.Lock()
				if u.TimeUsed > 10 {
					timeUp <- true
				}
				u.TimeUsed++
				mu.Unlock()
				time.Sleep(1 * time.Second)
			}
		}()
		process()
		done <- true
	}()
	select {

	case <-done:
		result = true
	case <-timeUp:
		log.Printf("times up")
		result = false
	case <-time.After(10 * time.Second):
		log.Printf("process killed")
		result = false
	}
	return result
}

func main() {

	start := time.Now()
	RunMockServer()
	fmt.Printf("Process took %s\n", time.Since(start))
}
